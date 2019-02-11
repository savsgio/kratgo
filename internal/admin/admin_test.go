package admin

import (
	"io"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/savsgio/kratgo/internal/cache"
	"github.com/savsgio/kratgo/internal/config"
	"github.com/savsgio/kratgo/internal/invalidator"

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

type Path struct {
	method string
	url    string
	view   atreugo.View
}

type mockServer struct {
	listenAndServeCalled bool
	logOutput            io.Writer

	paths []Path

	mu sync.RWMutex
}

type mockInvalidator struct {
	addCalled   bool
	startCalled bool

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

	return nil
}

func (mock *mockServer) ListenAndServe() error {
	mock.mu.Lock()
	mock.listenAndServeCalled = true
	mock.mu.Unlock()

	time.Sleep(250 * time.Millisecond)

	return nil
}

func (mock *mockServer) Path(httpMethod string, url string, viewFn atreugo.View) {
	mock.paths = append(mock.paths, Path{
		method: httpMethod,
		url:    url,
		view:   viewFn,
	})
}

func (mock *mockServer) SetLogOutput(output io.Writer) {
	mock.logOutput = output
}

func getPath(paths []Path, url, method string) *Path {
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
	return Config{
		FileConfig:  fileConfigAdmin(),
		Cache:       testCache,
		Invalidator: nil,
		HTTPScheme:  "http",
		LogLevel:    "error",
		LogOutput:   os.Stderr,
	}
}

func TestAdmin_init(t *testing.T) {
	serverMock := new(mockServer)

	admin := new(Admin)
	admin.server = serverMock
	admin.init()

	expectedPaths := []Path{
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
		p := getPath(expectedPaths, path.url, path.method)
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
