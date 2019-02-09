package proxy

import (
	"os"
	"sync"

	"github.com/savsgio/kratgo/internal/cache"
	"github.com/savsgio/kratgo/internal/invalidator"
	"github.com/savsgio/kratgo/internal/proxy/config"

	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/govaluate/v3"
	"github.com/valyala/fasthttp"
)

// Proxy ...
type Proxy struct {
	server      *fasthttp.Server
	cache       *cache.Cache
	invalidator *invalidator.Invalidator

	backends       []fetcher
	totalBackends  int
	currentBackend int

	httpScheme string

	nocacheRules []rule
	headersRules []headerRule

	log     *logger.Logger
	logFile *os.File

	tools sync.Pool

	cfg config.Config

	mu sync.RWMutex
}

type proxyTools struct {
	httpClient *httpClient
	params     *evalParams
	entry      *cache.Entry
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

type headerRule struct {
	rule

	action string
	name   string
	value  headerValue
}

// ###### INTERFACES ######

type fetcher interface {
	Do(req *fasthttp.Request, resp *fasthttp.Response) error
}
