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
	if e.Host == "" && e.Path == "" && e.Header.Key == "" {
		return invTypeInvalid
	}

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

func (i *Invalidator) invalidate(invalidationType, key string, entry *cache.Entry, e Entry) error {
	switch invalidationType {
	case invTypePath:
		return i.invalidateByPath(key, entry, e)
	case invTypeHeader:
		return i.invalidateByHeader(key, entry, e)
	case invTypePathHeader:
		return i.invalidateByPathHeader(key, entry, e)
	default:
		return i.invalidateByHost(key)
	}
}

func (i *Invalidator) invalidateAll(invalidationType string, e Entry) {
	atomic.AddInt32(&i.activeWorkers, 1)
	defer atomic.AddInt32(&i.activeWorkers, -1)

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

		i.invalidate(invalidationType, v.Key(), entry, e)

		entry.Reset()
	}

	cache.ReleaseEntry(entry)
}

func (i *Invalidator) invalidateHost(invalidationType string, e Entry) {
	atomic.AddInt32(&i.activeWorkers, 1)
	defer atomic.AddInt32(&i.activeWorkers, -1)

	key := e.Host
	entry := cache.AcquireEntry()

	err := i.cache.Get(key, entry)
	if err != nil {
		i.log.Errorf("Could not get responses from cache by key '%s': %v", key, err)
	} else if entry.Len() == 0 {
		return
	}

	i.invalidate(invalidationType, key, entry, e)

	cache.ReleaseEntry(entry)
}

func (i *Invalidator) waitAvailableWorkers() {
	for atomic.LoadInt32(&i.activeWorkers) > i.fileConfig.MaxWorkers {
		time.Sleep(100 * time.Millisecond)
	}
}

// Add ..
func (i *Invalidator) Add(e Entry) error {
	if t := i.invalidationType(e); t == invTypeInvalid {
		return errEmptyFields
	}

	i.chEntries <- e

	return nil
}

// Start ...
func (i *Invalidator) Start() {
	for e := range i.chEntries {
		invalidationType := i.invalidationType(e)

		i.waitAvailableWorkers()

		if e.Host != "" {
			go i.invalidateHost(invalidationType, e)
		} else {
			go i.invalidateAll(invalidationType, e)
		}

	}
}
