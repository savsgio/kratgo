package main

import (
	"flag"

	"github.com/savsgio/kratgo/internal/config"
	"github.com/savsgio/kratgo/internal/kratgo"
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

	kratgo, err := kratgo.New(*cfg)
	if err != nil {
		panic(err)
	}

	if err = kratgo.ListenAndServe(); err != nil {
		panic(err)
	}
}
