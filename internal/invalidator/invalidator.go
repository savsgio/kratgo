package invalidator

import (
	"sync/atomic"
	"time"

	"github.com/savsgio/kratgo/internal/cache"

	"github.com/allegro/bigcache"
	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/gotils"
)

// New ...
func New(cfg Config) *Invalidator {
	log := logger.New("kratgo-invalidator", cfg.LogLevel, cfg.LogOutput)

	i := &Invalidator{
		fileConfig: cfg.FileConfig,
		cache:      cfg.Cache,
		chEntries:  make(chan Entry),
		log:        log,
	}

	return i
}

func (i *Invalidator) invalidateByHost(cacheKey string) error {
	if err := i.cache.Del(cacheKey); err != nil && err != bigcache.ErrEntryNotFound {
		return err
	}

	return nil
}

func (i *Invalidator) invalidateByPath(cacheKey, path string, entry *cache.Entry) error {
	for _, resp := range entry.Responses {
		if path == gotils.B2S(resp.Path) {
			entry.DelResponse(resp.Path)
			return i.cache.Set(cacheKey, entry)
		}
	}

	return nil
}

func (i *Invalidator) invalidateByHeader(cacheKey, headerKey, headerValue string, entry *cache.Entry) error {
	for _, resp := range entry.Responses {
		for _, h := range resp.Headers {
			if gotils.B2S(h.Key) != headerKey && gotils.B2S(h.Value) == headerValue {
				entry.DelResponse(resp.Path)
				break
			}
		}
	}

	return i.cache.Set(cacheKey, entry)
}

func (i *Invalidator) waitAvailableWorkers() {
	for atomic.LoadInt32(&i.activeWorkers) > i.fileConfig.MaxWorkers {
		time.Sleep(100 * time.Millisecond)
	}
}

func (i *Invalidator) invalidate(e Entry) {
	atomic.AddInt32(&i.activeWorkers, 1)
	defer atomic.AddInt32(&i.activeWorkers, -1)

	if e.Host != "" {
		if err := i.invalidateByHost(e.Host); err != nil {
			i.log.Errorf("Could not invalidate cache by host '%s': %v", e.Host, err)
		}
		return
	}

	entry := cache.AcquireEntry()
	iter := i.cache.Iterator()

	for iter.SetNext() {
		v, err := iter.Value()
		if err != nil {
			i.log.Errorf("Could not get value from iterator: %v", err)
			continue
		}

		if err = cache.Unmarshal(entry, v.Value()); err != nil {
			i.log.Errorf("Could not decode cache value: %v", err)
			continue
		}

		k := v.Key()

		if e.Path != "" {
			if err = i.invalidateByPath(k, e.Path, entry); err != nil {
				i.log.Errorf("Could not invalidate cache by path '%s': %v", e.Path, err)
			}

		} else if e.Header.Key != "" {
			if err = i.invalidateByHeader(k, e.Header.Key, e.Header.Value, entry); err != nil {
				i.log.Errorf("Could not invalidate cache by header '%s = %s': %v", e.Header.Key, e.Header.Value, err)
			}
		}

		entry.Reset()
	}

	cache.ReleaseEntry(entry)
}

// Add ..
func (i *Invalidator) Add(e Entry) {
	i.chEntries <- e
}

// Start ...
func (i *Invalidator) Start() {
	for e := range i.chEntries {
		i.waitAvailableWorkers()
		go i.invalidate(e)
	}
}
