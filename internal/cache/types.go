package cache

import (
	"io"

	"github.com/savsgio/kratgo/internal/config"

	"github.com/allegro/bigcache"
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
