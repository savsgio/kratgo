package admin

import (
	logger "github.com/savsgio/go-logger"
	"github.com/valyala/fasthttp"
)

// New ...
func New(cfg Config) *Admin {
	a := new(Admin)
	a.fileConfig = cfg.FileConfig

	log := logger.New("kratgo-admin", cfg.LogLevel, cfg.LogOutput)

	a.server = &fasthttp.Server{
		Handler: a.httpHandler,
		Name:    "Kratgo",
		Logger:  log,
	}

	a.httpScheme = cfg.HTTPScheme

	a.cache = cfg.Cache
	a.invalidator = cfg.Invalidator
	a.log = log

	return a
}

// ListenAndServe ...
func (a *Admin) ListenAndServe() error {
	go a.invalidator.Start()

	a.log.Infof("Listening on: %s://%s/", a.httpScheme, a.fileConfig.Addr)

	return a.server.ListenAndServe(a.fileConfig.Addr)
}
