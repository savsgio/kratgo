package invalidator

import (
	"testing"

	"github.com/savsgio/kratgo/internal/cache"
)

func TestInvalidator_invalidateByHost(t *testing.T) {
	i := New(testConfig())

	key := "www.kratgo.com"

	resp := cache.AcquireResponse()
	resp.Path = []byte("/fast")

	entry := cache.AcquireEntry()
	entry.SetResponse(*resp)

	i.cache.Set(key, entry)

	i.invalidateByHost(key)

	entry.Reset()

	if err := i.cache.Get(key, entry); err != nil {
		t.Fatal(err)
	}

	if entry.Len() > 0 {
		t.Error("The cache has not been invalidate by host")
	}
}

func TestInvalidator_invalidateByPath(t *testing.T) {
	i := New(testConfig())

	path := "/fast"

	e := Entry{
		Path: path,
	}

	key := "www.kratgo.com"

	resp := cache.AcquireResponse()
	resp.Path = []byte(path)

	cacheEntry := cache.AcquireEntry()
	cacheEntry.SetResponse(*resp)

	i.cache.Set(key, cacheEntry)

	if err := i.invalidateByPath(key, cacheEntry, e); err != nil {
		t.Fatal(err)
	}

	cacheEntry.Reset()

	if err := i.cache.Get(key, cacheEntry); err != nil {
		t.Fatal(err)
	}

	if cacheEntry.HasResponse(resp.Path) {
		t.Error("The cache has not been invalidate by path")
	}
}

func TestInvalidator_invalidateByHeader(t *testing.T) {
	i := New(testConfig())

	path := "/fast"
	headerKey := []byte("kratgo")
	headerValue := []byte("fast")

	e := Entry{
		Path: path,
		Header: Header{
			Key:   string(headerKey),
			Value: string(headerValue),
		},
	}

	key := "www.kratgo.com"

	resp := cache.AcquireResponse()
	resp.SetHeader(headerKey, headerValue)

	cacheEntry := cache.AcquireEntry()
	cacheEntry.SetResponse(*resp)

	i.cache.Set(key, cacheEntry)

	if err := i.invalidateByHeader(key, cacheEntry, e); err != nil {
		t.Fatal(err)
	}

	cacheEntry.Reset()

	if err := i.cache.Get(key, cacheEntry); err != nil {
		t.Fatal(err)
	}

	if cacheEntry.HasResponse(resp.Path) {
		t.Error("The cache has not been invalidate by header")
	}
}

func TestInvalidator_invalidateByPathHeader(t *testing.T) {
	i := New(testConfig())

	path := "/fast"
	headerKey := []byte("kratgo")
	headerValue := []byte("fast")

	e := Entry{
		Path: path,
		Header: Header{
			Key:   string(headerKey),
			Value: string(headerValue),
		},
	}

	key := "www.kratgo.com"

	resp := cache.AcquireResponse()
	resp.Path = []byte(path)
	resp.SetHeader(headerKey, headerValue)

	cacheEntry := cache.AcquireEntry()
	cacheEntry.SetResponse(*resp)

	i.cache.Set(key, cacheEntry)

	if err := i.invalidateByPathHeader(key, cacheEntry, e); err != nil {
		t.Fatal(err)
	}

	cacheEntry.Reset()

	if err := i.cache.Get(key, cacheEntry); err != nil {
		t.Fatal(err)
	}

	if cacheEntry.HasResponse(resp.Path) {
		t.Error("The cache has not been invalidate by path and header")
	}
}
