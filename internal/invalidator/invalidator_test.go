package invalidator

import (
	"os"
	"testing"

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
		t string
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
					Header: Header{
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
					Header: Header{
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
				t.Errorf("Invalidator.invalidationType() = %s, want %s", got, tt.want.t)
			}
		})
	}
}
