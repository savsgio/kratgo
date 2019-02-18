package invalidator

import (
	"os"
	"testing"
	"time"

	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/kratgo/internal/cache"
	"github.com/savsgio/kratgo/internal/config"
)

var testCache *cache.Cache

func init() {
	c, err := cache.New(cache.Config{
		FileConfig: fileConfigCache(),
		LogLevel:   logger.ERROR,
		LogOutput:  os.Stderr,
	})
	if err != nil {
		panic(err)
	}

	testCache = c
}

func fileConfigInvalidator() config.Invalidator {
	return config.Invalidator{
		MaxWorkers: 1,
	}
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

func testConfig() Config {
	testCache.Reset()

	return Config{
		FileConfig: fileConfigInvalidator(),
		Cache:      testCache,
		LogLevel:   logger.ERROR,
		LogOutput:  os.Stderr,
	}
}

func TestInvalidator_invalidationType(t *testing.T) {
	type args struct {
		e Entry
	}
	type want struct {
		t invType
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Host",
			args: args{
				e: Entry{
					Host: "www.kratgo.com",
				},
			},
			want: want{
				t: invTypeHost,
			},
		},
		{
			name: "Path",
			args: args{
				e: Entry{
					Path: "/fast",
				},
			},
			want: want{
				t: invTypePath,
			},
		},
		{
			name: "Header",
			args: args{
				e: Entry{
					Header: EntryHeader{
						Key:   "X-Data",
						Value: "Fast",
					},
				},
			},
			want: want{
				t: invTypeHeader,
			},
		},
		{
			name: "PathHeader",
			args: args{
				e: Entry{
					Path: "/lightweight",
					Header: EntryHeader{
						Key:   "X-Data",
						Value: "Fast",
					},
				},
			},
			want: want{
				t: invTypePathHeader,
			},
		},
		{
			name: "Invalid",
			args: args{
				e: Entry{},
			},
			want: want{
				t: invTypeInvalid,
			},
		},
	}

	i := New(testConfig())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := i.invalidationType(tt.args.e); got != tt.want.t {
				t.Errorf("Invalidator.invalidationType() = %d, want %d", got, tt.want.t)
			}
		})
	}
}

func TestInvalidator_invalidate(t *testing.T) {
	type args struct {
		invalidationType invType
		entry            cache.Entry
		e                Entry
	}

	type want struct {
		foundInCache bool
	}

	host := "www.kratgo.com"
	responses := []cache.Response{
		{
			Path: []byte("/fast"),
			Body: []byte("Kratgo is not slow"),
			Headers: []cache.ResponseHeader{
				{Key: []byte("X-Data"), Value: []byte("1")},
			},
		},
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Host",
			args: args{
				invalidationType: invTypeHost,
				entry: cache.Entry{
					Responses: responses,
				},
				e: Entry{
					Host: host,
				},
			},
			want: want{
				foundInCache: false,
			},
		},
		{
			name: "Path",
			args: args{
				invalidationType: invTypePath,
				entry: cache.Entry{
					Responses: responses,
				},
				e: Entry{
					Host: host,
					Path: "/fast",
				},
			},
			want: want{
				foundInCache: false,
			},
		},
		{
			name: "PathHeader",
			args: args{
				invalidationType: invTypePathHeader,
				entry: cache.Entry{
					Responses: responses,
				},
				e: Entry{
					Host: host,
					Path: "/fast",
					Header: EntryHeader{
						Key:   "X-Data",
						Value: "1",
					},
				},
			},
			want: want{
				foundInCache: false,
			},
		},
		{
			name: "Header",
			args: args{
				invalidationType: invTypeHeader,
				entry: cache.Entry{
					Responses: responses,
				},
				e: Entry{
					Host: host,
					Header: EntryHeader{
						Key:   "X-Data",
						Value: "1",
					},
				},
			},
			want: want{
				foundInCache: false,
			},
		},
		{
			name: "Invalid",
			args: args{
				invalidationType: invTypeInvalid,
				entry: cache.Entry{
					Responses: responses,
				},
				e: Entry{
					Host: "www.fake.com",
				},
			},
			want: want{
				foundInCache: true,
			},
		},
	}

	i := New(testConfig())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := tt.args.e.Host

			if err := i.cache.Set(key, tt.args.entry); err != nil {
				t.Fatal(err)
			}

			if err := i.invalidate(tt.args.invalidationType, key, tt.args.entry, tt.args.e); err != nil {
				t.Fatal(err)
			}

			cacheEntry := cache.AcquireEntry()
			if err := i.cache.Get(key, cacheEntry); err != nil {
				t.Fatal(err)
			}

			length := cacheEntry.Len()
			if length > 0 && !tt.want.foundInCache {
				t.Errorf("Invalidator.invalidate() cache has not been invalidate, type = '%d'", tt.args.invalidationType)
			} else if length == 0 && tt.want.foundInCache {
				t.Errorf("Invalidator.invalidate() cache has been invalidate, type = '%d'", tt.args.invalidationType)
			}
		})
	}
}

func TestInvalidator_invalidateAll(t *testing.T) {
	type cacheData struct {
		host      string
		responses []cache.Response
	}

	path := []byte("/fast")

	host1 := "www.kratgo.com"
	responses1 := []cache.Response{
		{
			Path: path,
			Body: []byte("Kratgo is not slow"),
			Headers: []cache.ResponseHeader{
				{Key: []byte("X-Data"), Value: []byte("1")},
			},
		},
	}

	host2 := "www.cache-fast.com"
	responses2 := []cache.Response{
		{
			Path: path,
			Body: []byte("Kratgo is not slow"),
			Headers: []cache.ResponseHeader{
				{Key: []byte("X-Data"), Value: []byte("1")},
			},
		},
	}

	host3 := "www.high-performance.com"
	responses3 := []cache.Response{
		{
			Path: []byte("/kratgo"),
			Body: []byte("Kratgo is not slow"),
			Headers: []cache.ResponseHeader{
				{Key: []byte("X-Data"), Value: []byte("1")},
			},
		},
	}

	i := New(testConfig())

	i.cache.Set(host1, cache.Entry{Responses: responses1})
	i.cache.Set(host2, cache.Entry{Responses: responses2})
	i.cache.Set(host3, cache.Entry{Responses: responses3})

	i.invalidateAll(invTypePath, Entry{Path: string(path)})

	wantLength := 1
	length := i.cache.Len()
	if length != wantLength {
		t.Errorf("Invalidator.invalidateAll() cache length == '%d', want '%d'", length, wantLength)
	}
}

func TestInvalidator_invalidateHost(t *testing.T) {
	type cacheData struct {
		host      string
		responses []cache.Response
	}

	path := []byte("/fast")

	host1 := "www.kratgo.com"
	responses1 := []cache.Response{
		{
			Path: path,
			Body: []byte("Kratgo is not slow"),
			Headers: []cache.ResponseHeader{
				{Key: []byte("X-Data"), Value: []byte("1")},
			},
		},
	}

	host2 := "www.cache-fast.com"
	responses2 := []cache.Response{
		{
			Path: path,
			Body: []byte("Kratgo is not slow"),
			Headers: []cache.ResponseHeader{
				{Key: []byte("X-Data"), Value: []byte("1")},
			},
		},
	}

	i := New(testConfig())

	i.cache.Set(host1, cache.Entry{Responses: responses1})
	i.cache.Set(host2, cache.Entry{Responses: responses2})

	i.invalidateHost(invTypePath, Entry{Host: host1, Path: string(path)})

	wantLength := 1
	length := i.cache.Len()
	if length != wantLength {
		t.Errorf("Invalidator.invalidateAll() cache length == '%d', want '%d'", length, wantLength)
	}
}

// func TestInvalidator_waitAvailableWorkers(t *testing.T) {
// }

func TestInvalidator_Add(t *testing.T) {
	type args struct {
		entry Entry
	}
	type want struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Ok",
			args: args{
				entry: Entry{
					Host: "www.kratgo.com",
					Path: "/fast",
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "Error",
			args: args{
				entry: Entry{},
			},
			want: want{
				err: ErrEmptyFields,
			},
		},
	}

	i := New(testConfig())
	go i.Start()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := i.Add(tt.args.entry)

			if err != tt.want.err {
				t.Errorf("Invalidator.Add() error = %v, want %v", err, tt.want.err)
			}
		})
	}
}

func TestInvalidator_Start(t *testing.T) {
	key := "www.kratgo.com"
	path := "/fast"

	entry := Entry{
		Host: key,
		Path: path,
	}

	i := New(testConfig())
	i.chEntries = make(chan Entry, 1)

	cacheEntry := cache.AcquireEntry()
	cacheEntry.SetResponse(cache.Response{Path: []byte(path)})

	if err := i.cache.Set(key, *cacheEntry); err != nil {
		t.Fatal(err)
	}

	go i.Add(entry)

	go i.Start()

	time.Sleep(200 * time.Millisecond)

	if i.cache.Len() > 0 {
		t.Error("Invalidator.Start() invalidator has not been start")
	}
}
