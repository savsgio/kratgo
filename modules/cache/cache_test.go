package cache

import (
	"os"
	"reflect"
	"testing"
	"time"

	logger "github.com/savsgio/go-logger/v4"
	"github.com/savsgio/kratgo/modules/config"
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

	lifeWindoow := time.Duration(cfg.TTL) * time.Minute
	if bcConfig.LifeWindow != lifeWindoow {
		t.Errorf("bigcacheConfig() LifeWindow == '%d', want '%d'", bcConfig.LifeWindow, lifeWindoow)
	}

	cleanWindow := time.Duration(cfg.CleanFrequency) * time.Minute
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

func TestNew(t *testing.T) {
	type args struct {
		cfg Config
	}

	type want struct {
		err bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Ok",
			args: args{
				cfg: Config{
					FileConfig: config.Cache{
						TTL:              1,
						CleanFrequency:   1,
						MaxEntries:       1,
						MaxEntrySize:     1,
						HardMaxCacheSize: 10,
					},
					LogLevel:  logger.FATAL,
					LogOutput: os.Stderr,
				},
			},
			want: want{
				err: false,
			},
		},
		{
			name: "InvalidCleanFrequency",
			args: args{
				cfg: Config{
					FileConfig: config.Cache{
						TTL:              1,
						CleanFrequency:   0,
						MaxEntries:       1,
						MaxEntrySize:     1,
						HardMaxCacheSize: 10,
					},
					LogLevel:  logger.FATAL,
					LogOutput: os.Stderr,
				},
			},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.args.cfg)
			if (err != nil) != tt.want.err {
				t.Errorf("New() error = '%v', want '%v'", err, tt.want.err)
				return
			}

			if tt.want.err {
				return
			}

			if !reflect.DeepEqual(c.fileConfig, tt.args.cfg.FileConfig) {
				t.Errorf("New() fileConfig == '%v', want '%v'", c.fileConfig, tt.args.cfg.FileConfig)
			}

			if c.bc == nil {
				t.Errorf("New() bc is '%v'", nil)
			}
		})
	}
}

func TestCache_SetAndGetAndDel(t *testing.T) {
	e := getEntryTest()
	entry := AcquireEntry()

	k := "www.kratgo.com"

	err := testCache.Set(k, e)
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

	err := testCache.SetBytes(k, e)
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

	err := testCache.Set(k, e)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	iter := testCache.Iterator()
	if iter == nil {
		t.Errorf("Could not get iterator from cache")
	}
}

func TestCache_Len(t *testing.T) {
	e := getEntryTest()

	k := "www.kratgo.com"

	err := testCache.Set(k, e)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	wantLength := 1
	length := testCache.Len()
	if length != wantLength {
		t.Errorf("Cache.Len() == '%d', want '%d'", length, wantLength)
	}
}

func TestCache_Reset(t *testing.T) {
	e := getEntryTest()

	k := "www.kratgo.com"

	err := testCache.Set(k, e)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	testCache.Reset()

	wantLength := 0
	length := testCache.Len()
	if length != wantLength {
		t.Errorf("Cache.Len() == '%d', want '%d'", length, wantLength)
	}
}
