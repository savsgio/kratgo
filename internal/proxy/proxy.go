package proxy

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/savsgio/kratgo/internal/cache"
	"github.com/savsgio/kratgo/internal/invalidator"
	"github.com/savsgio/kratgo/internal/proxy/config"

	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/gotils"
	"github.com/savsgio/govaluate/v3"
	"github.com/valyala/fasthttp"
)

// go func() {
// 	time.Sleep(10 * time.Second)
// 	println("Invalidando la cache para crowne")
// 	invalidator.Add(&invalidator.Entry{
// 		Header: invalidator.Header{
// 			Key:   []byte("Theme"),
// 			Value: []byte("hans"),
// 		},
// 	})
// }()

// New ...
func New(cfg config.Config) (*Proxy, error) {
	p := new(Proxy)

	logOutput, err := getLogOutput(cfg.LogOutput)
	if err != nil {
		return nil, err
	}

	log := logger.New("kratgo", cfg.LogLevel, logOutput)

	s := &fasthttp.Server{
		Handler: p.handler,
		Name:    "Kratgo",
		Logger:  log,
	}

	hc := &fasthttp.HostClient{
		Addr: cfg.Proxy.BackendAddr,
	}

	c, err := cache.New(cache.Config{
		TTL:       cfg.Cache.TTL * time.Minute,
		LogLevel:  cfg.LogLevel,
		LogOutput: logOutput,
	})
	if err != nil {
		return nil, err
	}

	i := invalidator.New(invalidator.Config{
		Addr:       cfg.Invalidator.Addr,
		Cache:      c,
		MaxWorkers: cfg.Invalidator.MaxWorkers,
		LogLevel:   cfg.LogLevel,
		LogOutput:  logOutput,
	})

	p.server = s
	p.hostClient = hc
	p.cache = c
	p.invalidator = i
	p.httpScheme = "http"
	p.log = log
	p.cfg = cfg

	p.tools = sync.Pool{
		New: func() interface{} {
			return &proxyTools{
				fetcher: acquireFetcher(),
				params:  acquireEvalCtxParams(),
				entry:   cache.AcquireEntry(),
			}
		},
	}

	if err = p.parseNocacheRules(); err != nil {
		return nil, err
	}

	if err = p.parseHeadersRules(setHeaderAction, p.cfg.Proxy.Response.Headers.Set); err != nil {
		return nil, err
	}

	if err = p.parseHeadersRules(unsetHeaderAction, p.cfg.Proxy.Response.Headers.Unset); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Proxy) acquireTools() *proxyTools {
	return p.tools.Get().(*proxyTools)
}

func (p *Proxy) releaseTools(pt *proxyTools) {
	pt.fetcher.reset()
	pt.params.reset()
	pt.entry.Reset()

	p.tools.Put(pt)
}

func (p *Proxy) newEvaluableExpression(rule string) (*govaluate.EvaluableExpression, []ruleParam, error) {
	params := make([]ruleParam, 0)

	for config.ConfigVarRegex.MatchString(rule) {
		configKey, evalKey, evalSubKey := config.ParseConfigKeys(rule)

		rule = strings.Replace(rule, configKey, evalKey, -1)
		params = append(params, ruleParam{name: evalKey, subKey: evalSubKey})
	}

	expr, err := govaluate.NewEvaluableExpression(rule)
	return expr, params, err
}

func (p *Proxy) parseNocacheRules() error {
	for _, rule := range p.cfg.Proxy.Nocache {
		r := Rule{}

		expr, params, err := p.newEvaluableExpression(rule)
		if err != nil {
			return fmt.Errorf("Could not get the evaluable expression for rule '%s': %v", rule, err)
		}
		r.expr = expr
		r.params = append(r.params, params...)

		p.nocacheRules = append(p.nocacheRules, r)
	}

	return nil
}

func (p *Proxy) parseHeadersRules(action string, headers []config.Headers) error {
	for _, h := range headers {
		r := HeaderRule{action: action, name: h.Name}

		if h.When != "" {
			expr, params, err := p.newEvaluableExpression(h.When)
			if err != nil {
				return fmt.Errorf("Could not get the evaluable expression for rule '%s': %v", h.When, err)
			}
			r.expr = expr
			r.params = append(r.params, params...)
		}

		if action == setHeaderAction {
			_, evalKey, evalSubKey := config.ParseConfigKeys(h.Value)
			if evalKey != "" {
				r.value.value = evalKey
				r.value.subKey = evalSubKey
			} else {
				r.value.value = h.Value
			}
		}

		p.headersRules = append(p.headersRules, r)
	}

	return nil
}

func (p *Proxy) saveBackendResponse(cacheKey []byte, path string, resp *fasthttp.Response, entry *cache.Entry) error {
	r := entry.Response[path]
	r.Body = resp.Body()
	if r.Headers == nil {
		r.Headers = make(cache.ResponseHeaders)
	}
	resp.Header.VisitAll(func(k, v []byte) {
		r.Headers[gotils.B2S(k)] = v
	})

	entry.Response[path] = r

	if err := p.cache.SetBytes(cacheKey, entry); err != nil {
		return fmt.Errorf("Could not save response in cache for key '%s': %v", cacheKey, err)
	}

	return nil
}

func (p *Proxy) fetchFromBackend(cacheKey []byte, path string, ctx *fasthttp.RequestCtx, pt *proxyTools) error {
	if p.log.DebugEnabled() {
		p.log.Debugf("%s - %s", ctx.Method(), ctx.Path())
	}

	cloneHeaders(&pt.fetcher.req.Header, &ctx.Request.Header)
	pt.fetcher.req.Header.SetMethodBytes(ctx.Method())
	pt.fetcher.req.SetRequestURI(path)

	if err := pt.fetcher.Do(p.hostClient); err != nil {
		return fmt.Errorf("Could not fetch response from backend: %v", err)

	}

	pt.fetcher.processHeaderRules(p.headersRules, pt.params)
	pt.fetcher.req.Header.CopyTo(&ctx.Request.Header)
	pt.fetcher.resp.Header.CopyTo(&ctx.Response.Header)

	location := pt.fetcher.resp.Header.Peek(headerLocation)
	if len(location) > 0 {
		return nil
	}

	noCache, err := checkIfNoCache(pt.fetcher.req, pt.fetcher.resp, p.nocacheRules, pt.params)
	if err != nil {
		return err
	}

	ctx.SetStatusCode(pt.fetcher.resp.StatusCode())
	ctx.SetBody(pt.fetcher.resp.Body())

	if noCache || pt.fetcher.resp.StatusCode() != fasthttp.StatusOK {
		return nil
	}

	pt.fetcher.resp.Header.CopyTo(&ctx.Response.Header)

	return p.saveBackendResponse(cacheKey, path, &ctx.Response, pt.entry)
}

func (p *Proxy) handler(ctx *fasthttp.RequestCtx) {
	pt := p.acquireTools()

	path := gotils.B2S(ctx.URI().PathOriginal())
	cacheKey := ctx.Host()

	if noCache, err := checkIfNoCache(&ctx.Request, &ctx.Response, p.nocacheRules, pt.params); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		p.log.Error(err)

	} else if !noCache {
		if err := p.cache.GetBytes(cacheKey, pt.entry); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			p.log.Errorf("Could not get data from cache with key '%s': %v", cacheKey, err)

		} else if r, ok := pt.entry.Response[path]; ok {
			ctx.SetBody(r.Body)
			for k, v := range r.Headers {
				ctx.Response.Header.SetCanonical(gotils.S2B(k), v)
			}

			p.releaseTools(pt)
			return
		}
	}

	if err := p.fetchFromBackend(cacheKey, path, ctx, pt); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		p.log.Error(err)
	}

	p.releaseTools(pt)
}

// ListenAndServe ...
func (p *Proxy) ListenAndServe() error {
	defer p.logFile.Close()

	go p.invalidator.Start()

	p.log.Infof("Listening on: %s://%s/", p.httpScheme, p.cfg.Proxy.Addr)

	return p.server.ListenAndServe(p.cfg.Proxy.Addr)
}
