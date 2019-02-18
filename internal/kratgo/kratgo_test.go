package kratgo

import (
	"sync"
	"testing"
	"time"
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

func TestKratgo_Version(t *testing.T) {
	v := Version()
	if v != version {
		t.Errorf("Kratgo.Version() == '%s', want '%s'", v, version)
	}
}
