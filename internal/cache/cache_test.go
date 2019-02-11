package cache

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/savsgio/kratgo/internal/config"

	logger "github.com/savsgio/go-logger"
)

var testCache *Cache

func init() {
	c, err := New(Config{
		FileConfig: fileConfigCache(),
		LogLevel:   logger.ERROR,
		LogOutput:  os.Stderr,
	})
	if err != nil {
		panic(err)
	}

	testCache = c
}

func fileConfigCache() config.Cache {
	return config.Cache{
		TTL:              10,
		CleanFrequency:   5,
		MaxEntries:       5,
		MaxEntrySize:     20,
		HardMaxCacheSize: 30,
	}
}

func Test_bigcacheConfig(t *testing.T) {
	cfg := fileConfigCache()
	bcConfig := bigcacheConfig(cfg)

	if bcConfig.Shards != defaultBigcacheShards {
		t.Errorf("bigcacheConfig() Shards == '%d', want '%d'", bcConfig.Shards, defaultBigcacheShards)
	}

	lifeWindoow := cfg.TTL * time.Minute
	if bcConfig.LifeWindow != lifeWindoow {
		t.Errorf("bigcacheConfig() LifeWindow == '%d', want '%d'", bcConfig.LifeWindow, lifeWindoow)
	}

	cleanWindow := cfg.CleanFrequency * time.Minute
	if bcConfig.CleanWindow != cleanWindow {
		t.Errorf("bigcacheConfig() CleanWindow == '%d', want '%d'", bcConfig.CleanWindow, cleanWindow)
	}

	maxEntriesInWindow := cfg.MaxEntries
	if bcConfig.MaxEntriesInWindow != maxEntriesInWindow {
		t.Errorf("bigcacheConfig() MaxEntriesInWindow == '%d', want '%d'", bcConfig.MaxEntriesInWindow, maxEntriesInWindow)
	}

	maxEntriesSize := cfg.MaxEntrySize
	if bcConfig.MaxEntrySize != maxEntriesSize {
		t.Errorf("bigcacheConfig() MaxEntrySize == '%d', want '%d'", bcConfig.MaxEntrySize, maxEntriesSize)
	}

	verbose := false
	if bcConfig.Verbose != verbose {
		t.Errorf("bigcacheConfig() Verbose == '%v', want '%v'", bcConfig.Verbose, verbose)
	}

	hardMaxCacheSize := cfg.HardMaxCacheSize
	if bcConfig.HardMaxCacheSize != hardMaxCacheSize {
		t.Errorf("bigcacheConfig() HardMaxCacheSize == '%d', want '%d'", bcConfig.HardMaxCacheSize, hardMaxCacheSize)
	}
}

func TestCache_SetAndGetAndDel(t *testing.T) {
	e := getEntryTest()
	entry := AcquireEntry()

	k := "www.kratgo.com"

	err := testCache.Set(k, &e)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	err = testCache.Get(k, entry)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(e, *entry) {
		t.Errorf("The key '%s' has not been save in cache", k)
	}

	entry.Reset()

	err = testCache.Del(k)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	err = testCache.Get(k, entry)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if reflect.DeepEqual(e, *entry) {
		t.Errorf("The key '%s' has not been delete from cache", k)
	}
}

func TestCache_SetAndGetAndDel_Bytes(t *testing.T) {
	e := getEntryTest()
	entry := AcquireEntry()

	k := []byte("www.kratgo.com")

	err := testCache.SetBytes(k, &e)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	err = testCache.GetBytes(k, entry)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(e, *entry) {
		t.Errorf("The key '%s' has not been save in cache", k)
	}

	entry.Reset()

	err = testCache.DelBytes(k)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	err = testCache.GetBytes(k, entry)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if reflect.DeepEqual(e, *entry) {
		t.Errorf("The key '%s' has not been delete from cache", k)
	}
}

func TestCache_Iterator(t *testing.T) {
	e := getEntryTest()

	k := "www.kratgo.com"

	err := testCache.Set(k, &e)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	iter := testCache.Iterator()
	if iter == nil {
		t.Errorf("Could not get iterator from cache")
	}
}
