package config

// Config ...
type Config struct {
	Cache       Cache       `yaml:"cache"`
	Invalidator Invalidator `yaml:"invalidator"`
	Proxy       Proxy       `yaml:"proxy"`
	Admin       Admin       `yaml:"admin"`

	LogLevel  string `yaml:"logLevel"`
	LogOutput string `yaml:"logOutput"`
}

// Proxy ...
type Proxy struct {
	Addr         string        `yaml:"addr"`
	BackendAddrs []string      `yaml:"backendAddrs"`
	Response     ProxyResponse `yaml:"response"`
	Nocache      []string      `yaml:"nocache"`
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
	TTL              int `yaml:"ttl"`
	CleanFrequency   int `yaml:"cleanFrequency"`
	MaxEntries       int `yaml:"maxEntries"`
	MaxEntrySize     int `yaml:"maxEntrySize"`
	HardMaxCacheSize int `yaml:"hardMaxCacheSize"`
}

// Invalidator ...
type Invalidator struct {
	MaxWorkers int32 `yaml:"maxWorkers"`
}

// Admin ...
type Admin struct {
	Addr string `yaml:"addr"`
}
