package admin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/savsgio/atreugo/v7"
	logger "github.com/savsgio/go-logger"
)

// New ...
func New(cfg Config) (*Admin, error) {
	a := new(Admin)
	a.fileConfig = cfg.FileConfig

	logName := "kratgo-admin"

	log := logger.New(logName, cfg.LogLevel, cfg.LogOutput)

	addr := strings.Split(cfg.FileConfig.Addr, ":")
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		return nil, fmt.Errorf("Invalid address '%s': %v", cfg.FileConfig.Addr, err)
	}

	a.server = atreugo.New(&atreugo.Config{
		Host:    addr[0],
		Port:    port,
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
