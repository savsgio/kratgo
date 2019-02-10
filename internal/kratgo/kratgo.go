package kratgo

import (
	"fmt"

	"github.com/savsgio/kratgo/internal/admin"
	"github.com/savsgio/kratgo/internal/cache"
	"github.com/savsgio/kratgo/internal/config"
	"github.com/savsgio/kratgo/internal/invalidator"
	"github.com/savsgio/kratgo/internal/proxy"
)

// New ...
func New(cfg config.Config) (*Kratgo, error) {
	k := new(Kratgo)

	logFile, err := getLogOutput(cfg.LogOutput)
	if err != nil {
		return nil, err
	}
	k.logFile = logFile

	if cfg.Cache.CleanFrequency == 0 {
		return nil, fmt.Errorf("Cache.CleanFrequency configuration must be greater than 0")
	}

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

	i := invalidator.New(invalidator.Config{
		FileConfig: cfg.Invalidator,
		Cache:      c,
		LogLevel:   cfg.LogLevel,
		LogOutput:  logFile,
	})

	k.Admin = admin.New(admin.Config{
		FileConfig:  cfg.Admin,
		Cache:       c,
		Invalidator: i,
		HTTPScheme:  defaultHTTPScheme,
		LogLevel:    cfg.LogLevel,
		LogOutput:   logFile,
	})

	return k, nil
}

// ListenAndServe ...
func (k *Kratgo) ListenAndServe() error {
	defer k.logFile.Close()

	go k.Admin.ListenAndServe()

	return k.Proxy.ListenAndServe()
}
