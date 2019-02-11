package admin

import (
	"testing"

	"github.com/savsgio/atreugo/v7"
	"github.com/valyala/fasthttp"
)

func TestAdmin_invalidateView(t *testing.T) {
	invalidatorMock := new(mockInvalidator)

	admin, err := New(testConfig())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	admin.invalidator = invalidatorMock

	ctx := new(fasthttp.RequestCtx)
	actx := new(atreugo.RequestCtx)
	actx.RequestCtx = ctx

	actx.Request.Header.SetMethod("POST")
	actx.Request.SetBodyString("{\"action\": \"delete\", \"host\": \"www.crowne-demo.es\"}")

	err = admin.invalidateView(actx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !invalidatorMock.addCalled {
		t.Error("Admin.invalidateView() has not called to admin.invalidator.Add(...)")
	}

	expectedBody := "OK"
	respBody := string(actx.Response.Body())
	if respBody != expectedBody {
		t.Errorf("Admin.invalidateView() response body == '%s', want '%s'", respBody, expectedBody)
	}
}
