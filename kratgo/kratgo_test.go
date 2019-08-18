package kratgo

import (
	"sync"
	"testing"
	"time"

	logger "github.com/savsgio/go-logger"
	"github.com/savsgio/kratgo/modules/config"
)

type mockServer struct {
	listenAndServeCalled bool
	mu                   sync.RWMutex
}

func (mock *mockServer) ListenAndServe() error {
	mock.mu.Lock()
	mock.listenAndServeCalled = true
	mock.mu.Unlock()

	time.Sleep(250 * time.Millisecond)

	return nil
}

func TestKratgo_New(t *testing.T) {
	type args struct {
		cfg config.Config
	}

	type want struct {
		logFileName string
		err         bool
	}

	logFileName := "/tmp/test_kratgo.log"
	logLevel := logger.FATAL

	cfgAdmin := config.Admin{
		Addr: "localhost:9999",
	}
	cfgCache := config.Cache{
		TTL:              10,
		CleanFrequency:   5,
		MaxEntries:       5,
		MaxEntrySize:     20,
		HardMaxCacheSize: 30,
	}
	cfgInvalidator := config.Invalidator{
		MaxWorkers: 1,
	}
	cfgProxy := config.Proxy{
		Addr:         "localhost:8000",
		BackendAddrs: []string{"localhost:9990", "localhost:9991", "localhost:9993", "localhost:9994"},
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Ok",
			args: args{
				cfg: config.Config{
					Admin:       cfgAdmin,
					Cache:       cfgCache,
					Invalidator: cfgInvalidator,
					Proxy:       cfgProxy,
					LogLevel:    logLevel,
					LogOutput:   logFileName,
				},
			},
			want: want{
				logFileName: logFileName,
				err:         false,
			},
		},
		{
			name: "InvalidAdmin",
			args: args{
				cfg: config.Config{
					Cache:       cfgCache,
					Invalidator: cfgInvalidator,
					Proxy:       cfgProxy,
					LogLevel:    logLevel,
					LogOutput:   logFileName,
				},
			},
			want: want{
				err: true,
			},
		},
		{
			name: "InvalidCache",
			args: args{
				cfg: config.Config{
					Admin:       cfgAdmin,
					Invalidator: cfgInvalidator,
					Proxy:       cfgProxy,
					LogLevel:    logLevel,
					LogOutput:   logFileName,
				},
			},
			want: want{
				err: true,
			},
		},
		{
			name: "InvalidInvalidator",
			args: args{
				cfg: config.Config{
					Admin:     cfgAdmin,
					Cache:     cfgCache,
					Proxy:     cfgProxy,
					LogLevel:  logLevel,
					LogOutput: logFileName,
				},
			},
			want: want{
				err: true,
			},
		},
		{
			name: "InvalidProxy",
			args: args{
				cfg: config.Config{
					Admin:       cfgAdmin,
					Cache:       cfgCache,
					Invalidator: cfgInvalidator,
					LogLevel:    logLevel,
					LogOutput:   logFileName,
				},
			},
			want: want{
				err: true,
			},
		},
		{
			name: "InvalidLogOutput",
			args: args{
				cfg: config.Config{
					Admin:       cfgAdmin,
					Cache:       cfgCache,
					Invalidator: cfgInvalidator,
					Proxy:       cfgProxy,
					LogLevel:    logLevel,
				},
			},
			want: want{
				err: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := New(tt.args.cfg)
			if (err != nil) != tt.want.err {
				t.Fatalf("New() error == '%v', want '%v'", err, tt.want.err)
			}

			if tt.want.err {
				return
			}

			logName := k.logFile.Name()
			if logName != tt.want.logFileName {
				t.Errorf("Kratgo.New() log file == '%s', want '%s'", logName, tt.want.logFileName)
			}

			if k.Admin == nil {
				t.Errorf("Kratgo.New() Admin is '%v'", nil)
			}
			if k.Proxy == nil {
				t.Errorf("Kratgo.New() Proxy is '%v'", nil)
			}
		})
	}
}

func TestKratgo_ListenAndServe(t *testing.T) {
	proxyMock := new(mockServer)
	adminMock := new(mockServer)

	k := new(Kratgo)
	k.Proxy = proxyMock
	k.Admin = adminMock

	k.ListenAndServe()

	// Sleep to wait the gorutine start
	time.Sleep(500 * time.Millisecond)

	proxyMock.mu.RLock()
	defer proxyMock.mu.RUnlock()
	if !proxyMock.listenAndServeCalled {
		t.Error("Kratgo.ListenAndServe() proxy server is not listening")
	}

	adminMock.mu.RLock()
	defer adminMock.mu.RUnlock()
	if !adminMock.listenAndServeCalled {
		t.Error("Kratgo.ListenAndServe() admin server is not listening")
	}
}
