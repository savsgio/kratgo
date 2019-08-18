package kratgo

import (
	"github.com/savsgio/kratgo/modules/admin"
	"github.com/savsgio/kratgo/modules/cache"
	"github.com/savsgio/kratgo/modules/config"
	"github.com/savsgio/kratgo/modules/invalidator"
	"github.com/savsgio/kratgo/modules/proxy"
)

// New ...
func New(cfg config.Config) (*Kratgo, error) {
	k := new(Kratgo)

	logFile, err := getLogOutput(cfg.LogOutput)
	if err != nil {
		return nil, err
	}
	k.logFile = logFile

	c, err := cache.New(cache.Config{
		FileConfig: cfg.Cache,
		LogLevel:   cfg.LogLevel,
		LogOutput:  logFile,
	})
	if err != nil {
		return nil, err
	}

	p, err := proxy.New(proxy.Config{
		FileConfig: cfg.Proxy,
		Cache:      c,
		HTTPScheme: defaultHTTPScheme,
		LogLevel:   cfg.LogLevel,
		LogOutput:  logFile,
	})
	if err != nil {
		return nil, err
	}
	k.Proxy = p

	i, err := invalidator.New(invalidator.Config{
		FileConfig: cfg.Invalidator,
		Cache:      c,
		LogLevel:   cfg.LogLevel,
		LogOutput:  logFile,
	})
	if err != nil {
		return nil, err
	}

	a, err := admin.New(admin.Config{
		FileConfig:  cfg.Admin,
		Cache:       c,
		Invalidator: i,
		HTTPScheme:  defaultHTTPScheme,
		LogLevel:    cfg.LogLevel,
		LogOutput:   logFile,
	})
	if err != nil {
		return nil, err
	}
	k.Admin = a

	return k, nil
}

// ListenAndServe ...
func (k *Kratgo) ListenAndServe() error {
	defer k.logFile.Close()

	err := make(chan error, 1)

	go func() {
		err <- k.Admin.ListenAndServe()
	}()

	go func() {
		err <- k.Proxy.ListenAndServe()
	}()

	return <-err
}
