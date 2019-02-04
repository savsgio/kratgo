package config

import "time"

// Config ...
type Config struct {
	Proxy       Proxy       `yaml:"proxy"`
	Cache       Cache       `yaml:"cache"`
	Invalidator Invalidator `yaml:"invalidator"`

	LogLevel  string `yaml:"logLevel"`
	LogOutput string `yaml:"logOutput"`
}

// Proxy ...
type Proxy struct {
	Addr        string        `yaml:"addr"`
	BackendAddr string        `yaml:"backendAddr"`
	Response    ProxyResponse `yaml:"response"`
	Nocache     []string      `yaml:"nocache"`
}

// ProxyResponse ...
type ProxyResponse struct {
	Headers ProxyResponseHeaders `yaml:"headers"`
}

// ProxyResponseHeaders ...
type ProxyResponseHeaders struct {
	Set   []Headers `yaml:"set"`
	Unset []Headers `yaml:"unset"`
}

// Headers ...
type Headers struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
	When  string `yaml:"if"`
}

// Cache ...
type Cache struct {
	TTL time.Duration `yaml:"ttl"`
}

// Invalidator ...
type Invalidator struct {
	Addr       string `yaml:"addr"`
	MaxWorkers int32  `yaml:"maxWorkers"`
}
