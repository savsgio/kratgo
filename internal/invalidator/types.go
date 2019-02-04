package invalidator

import (
	"io"

	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/kratgo/internal/cache"
	"github.com/valyala/fasthttp"
)

// Config ...
type Config struct {
	Addr string

	Cache      *cache.Cache
	MaxWorkers int32

	LogLevel  string
	LogOutput io.Writer
}

// Invalidator ...
type Invalidator struct {
	server *fasthttp.Server
	cache  *cache.Cache

	maxWorkers    int32
	activeWorkers int32

	serverHTTPScheme string

	chEntries chan Entry
	log       *logger.Logger
	cfg       Config
}

// Header ...
type Header struct {
	Key   string
	Value string
}

// Entry ...
type Entry struct {
	Action string
	Host   string
	Path   string
	Header Header
}
