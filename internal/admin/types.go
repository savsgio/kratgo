package admin

import (
	"io"

	"github.com/savsgio/atreugo/v7"
	"github.com/savsgio/kratgo/internal/config"
	"github.com/savsgio/kratgo/internal/invalidator"

	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/kratgo/internal/cache"
)

// Config ...
type Config struct {
	FileConfig  config.Admin
	Cache       *cache.Cache
	Invalidator Invalidator

	HTTPScheme string

	LogLevel  string
	LogOutput io.Writer
}

// Admin ...
type Admin struct {
	fileConfig config.Admin

	server      Server
	cache       *cache.Cache
	invalidator Invalidator

	httpScheme string

	log *logger.Logger
}

// ###### INTERFACES ######

// Invalidator ...
type Invalidator interface {
	Start()
	Add(e invalidator.Entry)
}

// Server ...
type Server interface {
	ListenAndServe() error
	Path(httpMethod string, url string, viewFn atreugo.View)
	SetLogOutput(output io.Writer)
}
