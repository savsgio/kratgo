package admin

import (
	"github.com/savsgio/atreugo/v11"
	logger "github.com/savsgio/go-logger/v2"
)

// New ...
func New(cfg Config) (*Admin, error) {
	a := new(Admin)
	a.fileConfig = cfg.FileConfig

	logName := "kratgo-admin"
	log := logger.New(logName, cfg.LogLevel, cfg.LogOutput)

	a.server = atreugo.New(atreugo.Config{
		Addr:    cfg.FileConfig.Addr,
		LogName: logName,
	})
	a.server.SetLogOutput(cfg.LogOutput)

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
