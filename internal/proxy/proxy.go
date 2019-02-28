package proxy

import (
	"fmt"
	"strings"
	"sync"

	"github.com/savsgio/kratgo/internal/cache"
	"github.com/savsgio/kratgo/internal/config"

	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/govaluate/v3"
	"github.com/valyala/fasthttp"
)

// New ...
func New(cfg Config) (*Proxy, error) {
	if len(cfg.FileConfig.BackendAddrs) == 0 {
		return nil, fmt.Errorf("Proxy.BackendAddrs configuration is mandatory")
	}

	p := new(Proxy)
	p.fileConfig = cfg.FileConfig

	log := logger.New("kratgo", cfg.LogLevel, cfg.LogOutput)

	p.server = &fasthttp.Server{
		Handler: p.handler,
		Name:    "Kratgo",
		Logger:  log,
	}

	p.cache = cfg.Cache
	p.httpScheme = cfg.HTTPScheme
	p.log = log

	p.backends = make([]fetcher, len(p.fileConfig.BackendAddrs))
	for i, addr := range p.fileConfig.BackendAddrs {
		p.backends[i] = &fasthttp.HostClient{
			Addr: addr,
		}
	}
	p.totalBackends = len(p.fileConfig.BackendAddrs)

	p.tools = sync.Pool{
		New: func() interface{} {
			return &proxyTools{
				params: acquireEvalParams(),
				entry:  cache.AcquireEntry(),
			}
		},
	}

	if err := p.parseNocacheRules(); err != nil {
		return nil, err
	}

	if err := p.parseHeadersRules(setHeaderAction, p.fileConfig.Response.Headers.Set); err != nil {
		return nil, err
	}

	if err := p.parseHeadersRules(unsetHeaderAction, p.fileConfig.Response.Headers.Unset); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Proxy) acquireTools() *proxyTools {
	return p.tools.Get().(*proxyTools)
}

func (p *Proxy) releaseTools(pt *proxyTools) {
	pt.params.reset()
	pt.entry.Reset()

	p.tools.Put(pt)
}

func (p *Proxy) getBackend() fetcher {
	if p.totalBackends == 1 {
		return p.backends[0]
	}

	p.mu.Lock()

	if p.currentBackend >= p.totalBackends-1 {
		p.currentBackend = 0
	} else {
		p.currentBackend++
	}

	backend := p.backends[p.currentBackend]

	p.mu.Unlock()

	return backend
}

func (p *Proxy) newEvaluableExpression(rule string) (*govaluate.EvaluableExpression, []ruleParam, error) {
	params := make([]ruleParam, 0)

	for config.ConfigVarRegex.MatchString(rule) {
		configKey, evalKey, evalSubKey := config.ParseConfigKeys(rule)
		if configKey == "" && evalKey == "" && evalSubKey == "" {
			return nil, nil, fmt.Errorf("Invalid condition: %s", rule)
		}

		rule = strings.Replace(rule, configKey, evalKey, -1)
		params = append(params, ruleParam{name: evalKey, subKey: evalSubKey})
	}

	expr, err := govaluate.NewEvaluableExpression(rule)
	return expr, params, err
}

func (p *Proxy) parseNocacheRules() error {
	for _, ncRule := range p.fileConfig.Nocache {
		r := rule{}

		expr, params, err := p.newEvaluableExpression(ncRule)
		if err != nil {
			return fmt.Errorf("Could not get the evaluable expression for rule '%s': %v", ncRule, err)
		}
		r.expr = expr
		r.params = append(r.params, params...)

		p.nocacheRules = append(p.nocacheRules, r)
	}

	return nil
}

func (p *Proxy) parseHeadersRules(action typeHeaderAction, headers []config.Header) error {
	for _, h := range headers {
		r := headerRule{action: action, name: h.Name}

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

func (p *Proxy) saveBackendResponse(cacheKey, path []byte, resp *fasthttp.Response, entry *cache.Entry) error {
	r := cache.AcquireResponse()

	r.Path = append(r.Path, path...)
	r.Body = append(r.Body, resp.Body()...)
	resp.Header.VisitAll(func(k, v []byte) {
		r.SetHeader(k, v)
	})

	entry.SetResponse(*r)

	if err := p.cache.SetBytes(cacheKey, *entry); err != nil {
		return fmt.Errorf("Could not save response in cache for key '%s': %v", cacheKey, err)
	}

	cache.ReleaseResponse(r)

	return nil
}

func (p *Proxy) fetchFromBackend(cacheKey, path []byte, ctx *fasthttp.RequestCtx, pt *proxyTools) error {
	req, resp := &ctx.Request, &ctx.Response

	if p.log.DebugEnabled() {
		p.log.Debugf("%s - %s", ctx.Method(), ctx.Path())
	}

	req.Header.Set(proxyReqHeaderKey, proxyReqHeaderValue)
	for _, header := range hopHeaders {
		req.Header.Del(header)
	}

	if err := p.getBackend().Do(req, resp); err != nil {
		return fmt.Errorf("Could not fetch response from backend: %v", err)
	}

	if err := processHeaderRules(req, resp, p.headersRules, pt.params); err != nil {
		return fmt.Errorf("Could not process headers rules: %v", err)
	}

	location := ctx.Response.Header.Peek(headerLocation)
	if len(location) > 0 {
		return nil
	}

	noCache, err := checkIfNoCache(req, resp, p.nocacheRules, pt.params)
	if err != nil {
		return err
	}

	if noCache || ctx.Response.StatusCode() != fasthttp.StatusOK {
		return nil
	}

	return p.saveBackendResponse(cacheKey, path, &ctx.Response, pt.entry)
}

func (p *Proxy) handler(ctx *fasthttp.RequestCtx) {
	pt := p.acquireTools()

	path := ctx.URI().PathOriginal()
	cacheKey := ctx.Host()

	if noCache, err := checkIfNoCache(&ctx.Request, &ctx.Response, p.nocacheRules, pt.params); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		p.log.Error(err)

	} else if !noCache {
		if err := p.cache.GetBytes(cacheKey, pt.entry); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			p.log.Errorf("Could not get data from cache with key '%s': %v", cacheKey, err)

		} else if r := pt.entry.GetResponse(path); r != nil {
			ctx.SetBody(r.Body)
			for _, h := range r.Headers {
				ctx.Response.Header.SetCanonical(h.Key, h.Value)
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
	p.log.Infof("Listening on: %s://%s/", p.httpScheme, p.fileConfig.Addr)

	return p.server.ListenAndServe(p.fileConfig.Addr)
}
