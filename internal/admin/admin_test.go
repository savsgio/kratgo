package admin

import (
	"os"

	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/kratgo/internal/cache"
	"github.com/savsgio/kratgo/internal/config"
	"github.com/savsgio/kratgo/internal/invalidator"
)

var testCache *cache.Cache

func init() {
	c, err := cache.New(cache.Config{
		FileConfig: fileConfigCache(),
		LogLevel:   logger.ERROR,
		LogOutput:  os.Stderr,
	})
	if err != nil {
		panic(err)
	}

	testCache = c
}

type mockInvalidator struct {
	addCalled   bool
	startCalled bool
}

func (mock *mockInvalidator) Start() {
	mock.startCalled = true
}

func (mock *mockInvalidator) Add(e invalidator.Entry) {
	mock.addCalled = true
}

func fileConfigAdmin() config.Admin {
	return config.Admin{
		Addr: "localhost:9999",
	}
}

func fileConfigCache() config.Cache {
	return config.Cache{
		TTL:              10,
		CleanFrequency:   5,
		MaxEntries:       5,
		MaxEntrySize:     20,
		HardMaxCacheSize: 30,
	}
}

func testConfig(invalidatorMock Invalidator) Config {
	return Config{
		FileConfig:  fileConfigAdmin(),
		Cache:       testCache,
		Invalidator: invalidatorMock,
		HTTPScheme:  "http",
		LogLevel:    "error",
		LogOutput:   os.Stderr,
	}
}
