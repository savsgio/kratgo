package cache

import (
	"io"

	"github.com/savsgio/kratgo/modules/config"

	"github.com/allegro/bigcache/v2"
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
