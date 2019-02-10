package admin

import (
	"io"

	"github.com/savsgio/kratgo/internal/config"
	"github.com/savsgio/kratgo/internal/invalidator"

	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/kratgo/internal/cache"
	"github.com/valyala/fasthttp"
)

// Config ...
type Config struct {
	FileConfig  config.Admin
	Cache       *cache.Cache
	Invalidator *invalidator.Invalidator

	HTTPScheme string

	LogLevel  string
	LogOutput io.Writer
}

// Admin ...
type Admin struct {
	fileConfig config.Admin

	server      *fasthttp.Server
	cache       *cache.Cache
	invalidator *invalidator.Invalidator

	httpScheme string

	log *logger.Logger
}
