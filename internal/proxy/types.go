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
	hostClient  *fasthttp.HostClient
	cache       *cache.Cache
	invalidator *invalidator.Invalidator

	nocacheRules []Rule
	headersRules []HeaderRule

	httpScheme string

	log     *logger.Logger
	logFile *os.File

	tools sync.Pool

	cfg config.Config
}

type proxyTools struct {
	fetcher *fetcher
	params  *evalParams
	entry   *cache.Entry
}

type evalParams struct {
	p map[string]interface{}
}

type fetcher struct {
	req  *fasthttp.Request
	resp *fasthttp.Response

	executeHeaderRule bool
	redirectsCookies  [][]byte
}

type ruleParam struct {
	name   string
	subKey string
}

type headerValue struct {
	value  string
	subKey string
}

// Rule ...
type Rule struct {
	expr   *govaluate.EvaluableExpression
	params []ruleParam
}

// HeaderRule ...
type HeaderRule struct {
	action string
	name   string
	value  headerValue
	Rule
}
