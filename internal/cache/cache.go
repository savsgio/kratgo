package cache

import (
	"fmt"
	"time"

	"github.com/savsgio/kratgo/internal/config"

	"github.com/allegro/bigcache"
	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/gotils"
)

func bigcacheConfig(cfg config.Cache) bigcache.Config {
	return bigcache.Config{
		Shards:             defaultBigcacheShards,
		LifeWindow:         cfg.TTL * time.Minute,
		CleanWindow:        cfg.CleanFrequency * time.Minute,
		MaxEntriesInWindow: cfg.MaxEntries,
		MaxEntrySize:       cfg.MaxEntrySize,
		Verbose:            false,
		HardMaxCacheSize:   cfg.HardMaxCacheSize,
	}
}

// New ...
func New(cfg Config) (*Cache, error) {
	if cfg.FileConfig.CleanFrequency == 0 {
		return nil, fmt.Errorf("Cache.CleanFrequency configuration must be greater than 0")
	}

	c := new(Cache)
	c.fileConfig = cfg.FileConfig

	log := logger.New("kratgo-cache", cfg.LogLevel, cfg.LogOutput)

	bigcacheCFG := bigcacheConfig(c.fileConfig)
	bigcacheCFG.Logger = log
	bigcacheCFG.Verbose = cfg.LogLevel == logger.DEBUG

	c.bc, _ = bigcache.NewBigCache(bigcacheCFG)

	return c, nil
}

// Set ...
func (c *Cache) Set(key string, entry Entry) error {
	data, _ := Marshal(entry)

	return c.bc.Set(key, data)
}

// SetBytes ...
func (c *Cache) SetBytes(key []byte, entry Entry) error {
	return c.Set(gotils.B2S(key), entry)
}

// Get ...
func (c *Cache) Get(key string, dst *Entry) error {
	data, err := c.bc.Get(key)
	if err == bigcache.ErrEntryNotFound {
		return nil
	} else if err != nil {
		return err
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

// Len ...
func (c *Cache) Len() int {
	return c.bc.Len()
}

// Reset ...
func (c *Cache) Reset() error {
	return c.bc.Reset()
}
