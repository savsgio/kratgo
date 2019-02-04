package main

import (
	"flag"

	"github.com/savsgio/kratgo/internal/proxy"
	"github.com/savsgio/kratgo/internal/proxy/config"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "/etc/kratgo/kratgo.conf.yml", "Configuration file path")
	flag.Parse()
}

func main() {
	cfg, err := config.Parse(configFilePath)
	if err != nil {
		panic(err)
	}

	p, err := proxy.New(*cfg)
	if err != nil {
		panic(err)
	}

	if err = p.ListenAndServe(); err != nil {
		panic(err)
	}
}
