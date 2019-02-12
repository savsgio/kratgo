package invalidator

import (
	"sync/atomic"
	"time"

	"github.com/savsgio/kratgo/internal/cache"

	logger "github.com/savsgio/go-logger"
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

func (i *Invalidator) invalidationType(e Entry) string {
	if e.Path != "" {
		if e.Header.Key != "" {
			return invTypePathHeader
		}

		return invTypePath
	}

	if e.Header.Key != "" {
		return invTypeHeader
	}

	return invTypeHost
}

func (i *Invalidator) invalidateAll(e Entry) {
	atomic.AddInt32(&i.activeWorkers, 1)
	defer atomic.AddInt32(&i.activeWorkers, -1)

	invalidationType := i.invalidationType(e)

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

		key := v.Key()

		switch invalidationType {
		case invTypePath:
			if err = i.invalidateByPath(key, entry, e); err != nil {
				i.log.Errorf("Could not invalidate cache by path '%s': %v", e.Path, err)
			}
		case invTypeHeader:
			if err = i.invalidateByHeader(key, entry, e); err != nil {
				i.log.Errorf("Could not invalidate cache by header '%s = %s': %v", e.Header.Key, e.Header.Value, err)
			}
		case invTypePathHeader:
			if err = i.invalidateByPathHeader(key, entry, e); err != nil {
				i.log.Errorf("Could not invalidate cache by path '%s' and header '%s = %s': %v", e.Path, e.Header.Key, e.Header.Value, err)
			}
		}

		entry.Reset()
	}

	cache.ReleaseEntry(entry)
}

func (i *Invalidator) invalidate(e Entry) {
	atomic.AddInt32(&i.activeWorkers, 1)
	defer atomic.AddInt32(&i.activeWorkers, -1)

	invalidationType := i.invalidationType(e)

	key := e.Host
	entry := cache.AcquireEntry()

	err := i.cache.Get(key, entry)
	if err != nil {
		i.log.Errorf("Could not get responses from cache by key '%s': %v", key, err)
	} else if entry.Len() == 0 {
		return
	}

	switch invalidationType {
	case invTypePath:
		if err = i.invalidateByPath(key, entry, e); err != nil {
			i.log.Errorf("Could not invalidate cache by path '%s': %v", e.Path, err)
		}
	case invTypeHeader:
		if err = i.invalidateByHeader(key, entry, e); err != nil {
			i.log.Errorf("Could not invalidate cache by header '%s = %s': %v", e.Header.Key, e.Header.Value, err)
		}
	case invTypePathHeader:
		if err = i.invalidateByPathHeader(key, entry, e); err != nil {
			i.log.Errorf("Could not invalidate cache by path '%s' and header '%s = %s': %v", e.Path, e.Header.Key, e.Header.Value, err)
		}
	default:
		if err = i.invalidateByHost(key); err != nil {
			i.log.Errorf("Could not invalidate cache by host '%s': %v", key, err)
		}
	}

	cache.ReleaseEntry(entry)
}

func (i *Invalidator) waitAvailableWorkers() {
	for atomic.LoadInt32(&i.activeWorkers) > i.fileConfig.MaxWorkers {
		time.Sleep(100 * time.Millisecond)
	}
}

// Add ..
func (i *Invalidator) Add(e Entry) error {
	if e.Host == "" && e.Path == "" && e.Header.Key == "" {
		return errEmptyFields
	}

	i.chEntries <- e

	return nil
}

// Start ...
func (i *Invalidator) Start() {
	for e := range i.chEntries {
		i.waitAvailableWorkers()

		if e.Host != "" {
			go i.invalidate(e)
		} else {
			go i.invalidateAll(e)
		}

	}
}
