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
	Set   []Header `yaml:"set"`
	Unset []Header `yaml:"unset"`
}

// Header ...
type Header struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
	When  string `yaml:"if"`
}

// Cache ...
type Cache struct {
	TTL              time.Duration `yaml:"ttl"`
	CleanFrequency   time.Duration `yaml:"cleanFrequency"`
	MaxEntries       int           `yaml:"maxEntries"`
	MaxEntrySize     int           `yaml:"maxEntrySize"`
	HardMaxCacheSize int           `yaml:"hardMaxCacheSize"`
}

// Invalidator ...
type Invalidator struct {
	Addr       string `yaml:"addr"`
	MaxWorkers int32  `yaml:"maxWorkers"`
}
