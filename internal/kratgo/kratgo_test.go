package kratgo

import (
	"testing"
	"time"
)

type mockServer struct {
	listenAndServeCalled bool
}

func (mock *mockServer) ListenAndServe() error {
	mock.listenAndServeCalled = true

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

	if !proxyMock.listenAndServeCalled {
		t.Error("Admin.ListenAndServe() proxy server is not listening")
	}

	if !adminMock.listenAndServeCalled {
		t.Error("Admin.ListenAndServe() admin server is not listening")
	}
}
