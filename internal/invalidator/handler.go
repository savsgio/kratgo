package invalidator

import (
	"fmt"

	"github.com/savsgio/kratgo/internal/cache"

	"github.com/allegro/bigcache"
	"github.com/savsgio/gotils"
)

func (i *Invalidator) deleteCacheKey(cacheKey string) error {
	if err := i.cache.Del(cacheKey); err != nil && err != bigcache.ErrEntryNotFound {
		return err
	}

	return nil
}

func (i *Invalidator) invalidateByHost(cacheKey string) error {
	if err := i.deleteCacheKey(cacheKey); err != nil {
		return fmt.Errorf("Could not invalidate cache by host '%s': %v", cacheKey, err)
	}

	return nil
}

func (i *Invalidator) invalidateByPath(cacheKey string, cacheEntry cache.Entry, e Entry) error {
	path := gotils.S2B(e.Path)

	if !cacheEntry.HasResponse(path) {
		return nil
	}

	if cacheEntry.Len() == 1 {
		// Only delete the cache data for current key if remaining 1 response, to free memory
		if err := i.deleteCacheKey(cacheKey); err != nil {
			return fmt.Errorf("Could not invalidate cache by path '%s': %v", e.Path, err)
		}

		return nil
	}

	cacheEntry.DelResponse(path)

	if err := i.cache.Set(cacheKey, cacheEntry); err != nil {
		return fmt.Errorf("Could not invalidate cache by path '%s': %v", e.Path, err)
	}

	return nil
}

func (i *Invalidator) invalidateByHeader(cacheKey string, cacheEntry cache.Entry, e Entry) error {
	responses := cacheEntry.GetAllResponses()

	for _, resp := range responses {
		if !resp.HasHeader(gotils.S2B(e.Header.Key), gotils.S2B(e.Header.Value)) {
			continue
		}

		if cacheEntry.Len() == 1 {
			// Only delete the cache data for current key if remaining 1 response, to free memory
			if err := i.deleteCacheKey(cacheKey); err != nil {
				return fmt.Errorf("Could not invalidate cache by header '%s = %s': %v", e.Header.Key, e.Header.Value, err)
			}

			return nil
		}

		cacheEntry.DelResponse(resp.Path)
	}

	if err := i.cache.Set(cacheKey, cacheEntry); err != nil {
		return fmt.Errorf("Could not invalidate cache by header '%s = %s': %v", e.Header.Key, e.Header.Value, err)
	}

	return nil
}

func (i *Invalidator) invalidateByPathHeader(cacheKey string, cacheEntry cache.Entry, e Entry) error {
	path := gotils.S2B(e.Path)

	resp := cacheEntry.GetResponse(path)
	if resp == nil {
		return nil
	}

	if !resp.HasHeader(gotils.S2B(e.Header.Key), gotils.S2B(e.Header.Value)) {
		return nil
	}

	if cacheEntry.Len() == 1 {
		// Only delete the cache data for current key if remaining 1 response, to free memory
		if err := i.deleteCacheKey(cacheKey); err != nil {
			return fmt.Errorf("Could not invalidate cache by path '%s' and header '%s = %s': %v", e.Path, e.Header.Key, e.Header.Value, err)
		}

		return nil
	}

	cacheEntry.DelResponse(path)

	if err := i.cache.Set(cacheKey, cacheEntry); err != nil {
		return fmt.Errorf("Could not invalidate cache by path '%s' and header '%s = %s': %v", e.Path, e.Header.Key, e.Header.Value, err)
	}

	return nil
}
