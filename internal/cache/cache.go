package cache

import (
	"time"

	"github.com/allegro/bigcache"
	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/gotils"
)

func bigcacheConfig(config Config) bigcache.Config {
	return bigcache.Config{
		Shards:             1024,
		LifeWindow:         config.TTL,
		CleanWindow:        30 * time.Second,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            false,
		HardMaxCacheSize:   0,
		Logger:             logger.New("kratgo-cache", config.LogLevel, config.LogOutput),
	}
}

// New ...
func New(config Config) (*Cache, error) {
	bc, err := bigcache.NewBigCache(bigcacheConfig(config))
	if err != nil {
		return nil, err
	}

	return &Cache{bc: bc}, nil
}

// Set ...
func (c *Cache) Set(key string, entry *Entry) error {
	data, err := Marshal(entry)
	if err != nil {
		return err
	}

	return c.bc.Set(key, data)
}

// SetBytes ...
func (c *Cache) SetBytes(key []byte, entry *Entry) error {
	return c.Set(gotils.B2S(key), entry)
}

// Get ...
func (c *Cache) Get(key string, dst *Entry) error {
	data, err := c.bc.Get(key)
	if err != nil && err != bigcache.ErrEntryNotFound {
		return err
	} else if err == bigcache.ErrEntryNotFound {
		return nil
	}

	return Unmarshal(dst, data)
}

// GetBytes ...
func (c *Cache) GetBytes(key []byte, dst *Entry) error {
	return c.Get(gotils.B2S(key), dst)
}

// Del ...
func (c *Cache) Del(key string) error {
	return c.bc.Delete(key)
}

// DelBytes ...
func (c *Cache) DelBytes(key []byte) error {
	return c.Del(gotils.B2S(key))
}

// Iterator ...
func (c *Cache) Iterator() *bigcache.EntryInfoIterator {
	return c.bc.Iterator()
}
