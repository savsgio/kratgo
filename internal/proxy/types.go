package proxy

import (
	"io"
	"sync"

	"github.com/savsgio/kratgo/internal/cache"
	"github.com/savsgio/kratgo/internal/config"

	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/govaluate/v3"
	"github.com/valyala/fasthttp"
)

// Config ...
type Config struct {
	FileConfig config.Proxy
	Cache      *cache.Cache

	HTTPScheme string

	LogLevel  string
	LogOutput io.Writer
}

// Proxy ...
type Proxy struct {
	fileConfig config.Proxy

	server server
	cache  *cache.Cache

	backends       []fetcher
	totalBackends  int
	currentBackend int

	httpScheme string

	nocacheRules []rule
	headersRules []headerRule

	log   *logger.Logger
	tools sync.Pool
	mu    sync.RWMutex
}

type proxyTools struct {
	params *evalParams
	entry  *cache.Entry
}

type httpClient struct {
	req  *fasthttp.Request
	resp *fasthttp.Response

	executeHeaderRule bool
}

type evalParams struct {
	p map[string]interface{}
}

type ruleParam struct {
	name   string
	subKey string
}

type headerValue struct {
	value  string
	subKey string
}

type rule struct {
	expr   *govaluate.EvaluableExpression
	params []ruleParam
}

type typeHeaderAction int

type headerRule struct {
	rule

	action typeHeaderAction
	name   string
	value  headerValue
}

// ###### INTERFACES ######

type fetcher interface {
	Do(req *fasthttp.Request, resp *fasthttp.Response) error
}

// Server ...
type server interface {
	ListenAndServe(addr string) error
}
