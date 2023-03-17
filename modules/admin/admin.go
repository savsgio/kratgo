package admin

import (
	"github.com/savsgio/atreugo/v11"
	logger "github.com/savsgio/go-logger/v4"
)

// New ...
func New(cfg Config) (*Admin, error) {
	a := new(Admin)
	a.fileConfig = cfg.FileConfig

	log := logger.New(cfg.LogLevel, cfg.LogOutput, logger.Field{Key: "type", Value: "admin"})

	a.server = atreugo.New(atreugo.Config{
		Addr:   cfg.FileConfig.Addr,
		Logger: log,
	})

	a.httpScheme = cfg.HTTPScheme
	a.cache = cfg.Cache
	a.invalidator = cfg.Invalidator
	a.log = log

	a.init()

	return a, nil
}

func (a *Admin) init() {
	a.server.Path("POST", "/invalidate/", a.invalidateView)
}

// ListenAndServe ...
func (a *Admin) ListenAndServe() error {
	go a.invalidator.Start()

	return a.server.ListenAndServe()
}
