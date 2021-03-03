package cache

import (
	"io"

	"github.com/allegro/bigcache/v3"
	"github.com/savsgio/kratgo/modules/config"
)

// Config ...
type Config struct {
	FileConfig config.Cache

	LogLevel  string
	LogOutput io.Writer
}

// Cache ...
type Cache struct {
	fileConfig config.Cache

	bc *bigcache.BigCache
}
