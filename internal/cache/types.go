package cache

import (
	"io"
	"time"

	"github.com/allegro/bigcache"
)

// Config ...
type Config struct {
	TTL time.Duration

	LogLevel  string
	LogOutput io.Writer
}

// Cache ...
type Cache struct {
	bc *bigcache.BigCache
}
