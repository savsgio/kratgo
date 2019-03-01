package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/savsgio/kratgo/kratgo"
	"github.com/savsgio/kratgo/modules/config"
)

var configFilePath string

func init() {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Print Kratgo version")

	flag.StringVar(&configFilePath, "config", "/etc/kratgo/kratgo.conf.yml", "Configuration file path")

	flag.Parse()

	if showVersion {
		fmt.Printf("Kratgo version: %s\n", kratgo.Version())
		fmt.Printf("Go version: %s\n", runtime.Version())
		os.Exit(0)
	}
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
