package invalidator

import (
	"sync/atomic"
	"time"

	"github.com/savsgio/kratgo/internal/cache"

	"github.com/allegro/bigcache"
	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/gotils"
	"github.com/valyala/fasthttp"
)

// New ...
func New(config Config) *Invalidator {
	log := logger.New("kratgo-invalidator", config.LogLevel, config.LogOutput)

	i := &Invalidator{
		cache:            config.Cache,
		maxWorkers:       config.MaxWorkers,
		serverHTTPScheme: "http",
		chEntries:        make(chan Entry),
		log:              log,
		cfg:              config,
	}

	i.server = &fasthttp.Server{
		Handler: i.httpHandler,
		Name:    "Kratgo",
		Logger:  log,
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
	for entryPath := range entry.Response {
		if path == entryPath {
			delete(entry.Response, path)
			return i.cache.Set(cacheKey, entry)
		}
	}

	return nil
}

func (i *Invalidator) invalidateByHeader(cacheKey, headerKey, headerValue string, entry *cache.Entry) error {
	for path, response := range entry.Response {
		if v, ok := response.Headers[headerKey]; ok && gotils.B2S(v) == headerValue {
			delete(entry.Response, path)
		}
	}

	return i.cache.Set(cacheKey, entry)
}

func (i *Invalidator) waitAvailableWorkers() {
	for atomic.LoadInt32(&i.activeWorkers) > i.maxWorkers {
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
	go func() {
		if err := i.ListenAndServe(); err != nil {
			panic(err)
		}

	}()

	for e := range i.chEntries {
		i.waitAvailableWorkers()
		go i.invalidate(e)
	}
}

// ListenAndServe ...
func (i *Invalidator) ListenAndServe() error {
	i.log.Infof("Listening on: %s://%s/", i.serverHTTPScheme, i.cfg.Addr)

	return i.server.ListenAndServe(i.cfg.Addr)
}
