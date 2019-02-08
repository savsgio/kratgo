package cache

import (
	"io"
	"time"

	"github.com/allegro/bigcache"
)

// Config ...
type Config struct {
	// Time after which entry can be evicted
	TTL time.Duration

	// Interval between removing expired entries (clean up).
	// If set to <= 0 then no action is performed. Setting to < 1 second is counterproductive — cache has a one second resolution.
	CleanFrequency time.Duration

	// Max number of entries in cache. Used only to calculate initial size for cache shards.
	// When proper value is set then additional memory allocation does not occur.
	MaxEntries int

	// Max size of entry in bytes. Used only to calculate initial size for cache shards.
	MaxEntrySize int

	// HardMaxCacheSize is a limit for cache size in MB. Cache will not allocate more memory than this limit.
	// It can protect application from consuming all available memory on machine, therefore from running OOM Killer.
	// Default value is 0 which means unlimited size. When the limit is higher than 0 and reached then
	// the oldest entries are overridden for the new ones.
	HardMaxCacheSize int

	// Verbose mode prints information about new memory allocation
	Verbose bool

	LogLevel  string
	LogOutput io.Writer
}

// Cache ...
type Cache struct {
	bc *bigcache.BigCache
}
