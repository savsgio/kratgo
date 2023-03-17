package admin

import (
	"io"

	"github.com/savsgio/atreugo/v11"
	logger "github.com/savsgio/go-logger/v4"
	"github.com/savsgio/kratgo/modules/cache"
	"github.com/savsgio/kratgo/modules/config"
	"github.com/savsgio/kratgo/modules/invalidator"
)

// Config ...
type Config struct {
	FileConfig  config.Admin
	Cache       *cache.Cache
	Invalidator Invalidator

	HTTPScheme string

	LogLevel  logger.Level
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
	Add(e invalidator.Entry) error
}

// Server ...
type Server interface {
	ListenAndServe() error
	Path(httpMethod string, url string, viewFn atreugo.View) *atreugo.Path
}
