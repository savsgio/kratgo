package admin

import (
	"io"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/savsgio/kratgo/modules/cache"
	"github.com/savsgio/kratgo/modules/config"
	"github.com/savsgio/kratgo/modules/invalidator"

	"github.com/savsgio/atreugo/v7"
	logger "github.com/savsgio/go-logger"
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

type mockPath struct {
	method string
	url    string
	view   atreugo.View
}

type mockServer struct {
	listenAndServeCalled bool
	logOutput            io.Writer

	paths []mockPath

	mu sync.RWMutex
}

func (mock *mockServer) ListenAndServe() error {
	mock.mu.Lock()
	mock.listenAndServeCalled = true
	mock.mu.Unlock()

	time.Sleep(250 * time.Millisecond)

	return nil
}

func (mock *mockServer) Path(httpMethod string, url string, viewFn atreugo.View) {
	mock.paths = append(mock.paths, mockPath{
		method: httpMethod,
		url:    url,
		view:   viewFn,
	})
}

func (mock *mockServer) SetLogOutput(output io.Writer) {
	mock.logOutput = output
}

type mockInvalidator struct {
	addCalled   bool
	startCalled bool
	err         error

	mu sync.RWMutex
}

func (mock *mockInvalidator) Start() {
	mock.mu.Lock()
	mock.startCalled = true
	mock.mu.Unlock()
}

func (mock *mockInvalidator) Add(e invalidator.Entry) error {
	mock.mu.Lock()
	mock.addCalled = true
	mock.mu.Unlock()

	return mock.err
}

func getMockPath(paths []mockPath, url, method string) *mockPath {
	for _, v := range paths {
		if v.url == url && v.method == method {
			return &v
		}
	}

	return nil
}

func fileConfigAdmin() config.Admin {
	return config.Admin{
		Addr: "localhost:9999",
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
		FileConfig:  fileConfigAdmin(),
		Cache:       testCache,
		Invalidator: nil,
		HTTPScheme:  "http",
		LogLevel:    logger.FATAL,
		LogOutput:   os.Stderr,
	}
}

func TestAdmin_New(t *testing.T) {
	type args struct {
		cfg Config
	}

	type want struct {
		err bool
	}

	logLevel := logger.FATAL
	logOutput := os.Stderr
	httpScheme := "http"
	invalidatorMock := new(mockInvalidator)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Ok",
			args: args{
				cfg: Config{
					FileConfig: config.Admin{
						Addr: "localhost:9999",
					},
					Cache:       testCache,
					Invalidator: invalidatorMock,
					HTTPScheme:  httpScheme,
					LogLevel:    logLevel,
					LogOutput:   logOutput,
				},
			},
			want: want{
				err: false,
			},
		},
		{
			name: "InvalidAddress",
			args: args{
				cfg: Config{
					FileConfig: config.Admin{
						Addr: "localhost",
					},
					Cache:       testCache,
					Invalidator: invalidatorMock,
					HTTPScheme:  httpScheme,
					LogLevel:    logLevel,
					LogOutput:   logOutput,
				},
			},
			want: want{
				err: true,
			},
		},
		{
			name: "InvalidPort",
			args: args{
				cfg: Config{
					FileConfig: config.Admin{
						Addr: "localhost:",
					},
					Cache:       testCache,
					Invalidator: invalidatorMock,
					HTTPScheme:  httpScheme,
					LogLevel:    logLevel,
					LogOutput:   logOutput,
				},
			},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := New(tt.args.cfg)
			if (err != nil) != tt.want.err {
				t.Fatalf("New() error == '%v', want '%v'", err, tt.want.err)
			}

			if tt.want.err {
				return
			}

			if a.server == nil {
				t.Errorf("New() server is '%v'", nil)
			}

			if a.httpScheme != httpScheme {
				t.Errorf("New() httpScheme == '%s', want '%s'", a.httpScheme, httpScheme)
			}

			adminCachePtr := reflect.ValueOf(a.cache).Pointer()
			testCachePtr := reflect.ValueOf(testCache).Pointer()
			if adminCachePtr != testCachePtr {
				t.Errorf("New() cache == '%d', want '%d'", adminCachePtr, testCachePtr)
			}

			adminInvalidatorPtr := reflect.ValueOf(a.invalidator).Pointer()
			invalidatorPtr := reflect.ValueOf(invalidatorMock).Pointer()
			if adminInvalidatorPtr != invalidatorPtr {
				t.Errorf("New() invalidator == '%d', want '%d'", adminInvalidatorPtr, invalidatorPtr)
			}

			if a.log == nil {
				t.Errorf("New() log is '%v'", nil)
			}
		})
	}
}

func TestAdmin_init(t *testing.T) {
	serverMock := new(mockServer)

	admin := new(Admin)
	admin.server = serverMock
	admin.init()

	expectedPaths := []mockPath{
		{
			method: "POST",
			url:    "/invalidate/",
			view:   admin.invalidateView,
		},
	}

	if len(expectedPaths) != len(serverMock.paths) {
		t.Fatalf("Admin.server.init() registered paths == '%v', want '%v'", serverMock.paths, expectedPaths)
	}

	for _, path := range serverMock.paths {
		p := getMockPath(expectedPaths, path.url, path.method)
		if p == nil {
			t.Errorf("Admin.server.path() method == '%s', want '%s'", path.method, p.method)
			t.Errorf("Admin.server.path() url == '%s', want '%s'", path.url, p.url)

		} else {
			if reflect.ValueOf(path.view).Pointer() != reflect.ValueOf(p.view).Pointer() {
				t.Errorf("Admin.server.path() url == '%p', want '%p'", path.view, p.view)
			}
		}
	}

}

func TestAdmin_ListenAndServe(t *testing.T) {
	serverMock := new(mockServer)
	invalidatorMock := new(mockInvalidator)

	admin := new(Admin)
	admin.server = serverMock
	admin.invalidator = invalidatorMock

	admin.ListenAndServe()

	invalidatorMock.mu.RLock()
	defer invalidatorMock.mu.RUnlock()
	if !invalidatorMock.startCalled {
		t.Error("Admin.ListenAndServe() invalidator is not start")
	}

	serverMock.mu.RLock()
	defer serverMock.mu.RUnlock()
	if !serverMock.listenAndServeCalled {
		t.Error("Admin.ListenAndServe() server is not listening")
	}
}
