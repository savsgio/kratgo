package invalidator

import (
	"github.com/allegro/bigcache"
	"github.com/savsgio/gotils"
	"github.com/savsgio/kratgo/internal/cache"
)

func (i *Invalidator) invalidateByHost(cacheKey string) error {
	if err := i.cache.Del(cacheKey); err != nil && err != bigcache.ErrEntryNotFound {
		return err
	}

	return nil
}

func (i *Invalidator) invalidateByPath(cacheKey string, cacheEntry *cache.Entry, e Entry) error {
	path := gotils.S2B(e.Path)

	if !cacheEntry.HasResponse(path) {
		return nil
	}

	cacheEntry.DelResponse(path)

	return i.cache.Set(cacheKey, cacheEntry)
}

func (i *Invalidator) invalidateByHeader(cacheKey string, cacheEntry *cache.Entry, e Entry) error {
	responses := cacheEntry.GetAllResponses()

	for _, resp := range responses {
		if !resp.HasHeader(gotils.S2B(e.Header.Key), gotils.S2B(e.Header.Value)) {
			continue
		}

		cacheEntry.DelResponse(resp.Path)
	}

	return i.cache.Set(cacheKey, cacheEntry)
}

func (i *Invalidator) invalidateByPathHeader(cacheKey string, cacheEntry *cache.Entry, e Entry) error {
	path := gotils.S2B(e.Path)

	resp := cacheEntry.GetResponse(path)
	if resp == nil {
		return nil
	}

	if !resp.HasHeader(gotils.S2B(e.Header.Key), gotils.S2B(e.Header.Value)) {
		return nil
	}

	cacheEntry.DelResponse(path)

	return i.cache.Set(cacheKey, cacheEntry)
}
