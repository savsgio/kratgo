package proxy

import (
	"bytes"
	"testing"

	"github.com/savsgio/kratgo/internal/config"

	"github.com/savsgio/gotils"
	"github.com/valyala/fasthttp"
)

type mockHTTPClient struct {
	called bool

	body       []byte
	headers    map[string][]byte
	statusCode int
}

func (mock *mockHTTPClient) Do(req *fasthttp.Request, resp *fasthttp.Response) error {
	mock.called = true

	resp.SetBody(mock.body)
	resp.SetStatusCode(mock.statusCode)

	for k, v := range mock.headers {
		resp.Header.SetCanonical(gotils.S2B(k), v)
	}

	return nil
}

func TestHTTPClient_do(t *testing.T) {
	hc := acquireHTTPClient()
	mockFetcher := &mockHTTPClient{}

	hc.do(mockFetcher)

	if !mockFetcher.called {
		t.Error("httpClient.do() fetcher.Do() is not called")
	}
}

func TestHTTPClient_setMethodBytes(t *testing.T) {
	hc := acquireHTTPClient()
	method := []byte("POST")

	hc.setMethodBytes(method)

	reqMethod := hc.req.Header.Method()
	if !bytes.Equal(reqMethod, method) {
		t.Errorf("httpClient.setMethodBytes() method == '%s', want '%s'", reqMethod, method)
	}
}

func TestHTTPClient_setRequestURIBytes(t *testing.T) {
	hc := acquireHTTPClient()
	path := []byte("/kratgo")

	hc.setRequestURIBytes(path)

	reqPath := hc.req.URI().PathOriginal()
	if !bytes.Equal(reqPath, path) {
		t.Errorf("httpClient.setMethodBytes() path == '%s', want '%s'", reqPath, path)
	}
}

func TestHTTPClient_copyReqHeaderTo(t *testing.T) {
	hc := acquireHTTPClient()
	hc.req.Header.Set("Kratgo", "Cache")
	hc.req.Header.Set("HTTP", "Fast")

	req := fasthttp.AcquireRequest()

	hc.copyReqHeaderTo(&req.Header)

	hc.req.Header.VisitAll(func(k, v []byte) {
		vh := req.Header.PeekBytes(k)
		if !bytes.Equal(vh, v) {
			t.Errorf("httpClient.copyReqHeaderTo() header '%s' is not copied", k)
		}
	})
}

func TestHTTPClient_copyRespHeaderTo(t *testing.T) {
	hc := acquireHTTPClient()
	hc.req.Header.Set("Kratgo", "Cache")
	hc.req.Header.Set("HTTP", "Fast")

	resp := fasthttp.AcquireResponse()

	hc.copyRespHeaderTo(&resp.Header)

	hc.resp.Header.VisitAll(func(k, v []byte) {
		vh := resp.Header.PeekBytes(k)
		if !bytes.Equal(vh, v) {
			t.Errorf("httpClient.copyRespHeaderTo() header '%s' is not copied", k)
		}
	})
}

func TestHTTPClient_respHeaderPeek(t *testing.T) {
	k := "Kratgo"
	v := "Cache Fast"

	hc := acquireHTTPClient()
	hc.resp.Header.Set(k, v)

	vh := hc.respHeaderPeek(k)

	if string(vh) != v {
		t.Errorf("httpClient.respHeaderPeek() header value of '%s' == '%s', want '%s'", k, vh, v)
	}
}

func TestHTTPClient_statusCode(t *testing.T) {
	hc := acquireHTTPClient()
	statusCode := 404

	hc.resp.SetStatusCode(statusCode)

	respStatusCode := hc.statusCode()
	if respStatusCode != statusCode {
		t.Errorf("httpClient.statusCode() status code == '%d', want '%d'", respStatusCode, statusCode)
	}
}

func TestHTTPClient_body(t *testing.T) {
	hc := acquireHTTPClient()
	body := []byte("Test Kratgo")

	hc.resp.SetBody(body)

	respBody := hc.body()
	if !bytes.Equal(respBody, body) {
		t.Errorf("httpClient.body() body == '%s', want '%s'", respBody, body)
	}
}

func TestHTTPClient_processHeaderRules(t *testing.T) {
	setName1 := "Kratgo"
	setValue1 := "Fast"
	setWhen1 := "$(resp.header::Content-Type) == 'text/html'"

	setName2 := "X-Data"
	setValue2 := "1"

	// This header not fulfill the condition because $(req.header::X-Data) == '123'
	setName3 := "X-NotSet"
	setValue3 := "yes"
	setWhen3 := "$(req.header::X-Data) != '123'"
	// ----

	unsetName1 := "X-Delete"
	unsetWhen1 := "$(resp.header::Content-Type) == 'text/html'"

	unsetName2 := "X-MyHeader"

	setHeadersRulesConfig := []config.Header{
		{
			Name:  setName1,
			Value: setValue1,
			When:  setWhen1,
		},
		{
			Name:  setName2,
			Value: setValue2,
		},
		{
			Name:  setName3,
			Value: setValue3,
			When:  setWhen3,
		},
	}
	unsetHeadersRulesConfig := []config.Header{
		{
			Name: unsetName1,
			When: unsetWhen1,
		},
		{
			Name: unsetName2,
		},
	}

	p, _ := New(testConfig())
	p.parseHeadersRules(setHeaderAction, setHeadersRulesConfig)
	p.parseHeadersRules(unsetHeaderAction, unsetHeadersRulesConfig)

	params := acquireEvalParams()

	hc := acquireHTTPClient()

	hc.resp.Header.Set(unsetName1, "data")
	hc.resp.Header.Set("FakeHeader", "fake data")
	hc.resp.Header.Set("Content-Type", "text/html")
	hc.req.Header.Set("X-Data", "123")

	err := hc.processHeaderRules(p.headersRules, params)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if v := hc.resp.Header.Peek(setName1); string(v) != setValue1 {
		t.Errorf("httpClient.processHeaderRules() not set header '%s' with value '%s', want '%s==%s'",
			setName1, setValue1, setName1, v)
	}

	if v := hc.resp.Header.Peek(setName2); string(v) != setValue2 {
		t.Errorf("httpClient.processHeaderRules() not set header '%s' with value '%s', want '%s==%s'",
			setName2, setValue2, setName2, v)
	}

	if v := hc.resp.Header.Peek(setName3); len(v) > 0 {
		t.Errorf("httpClient.processHeaderRules() header '%s' is setted but not fulfill the condition", setName3)
	}

	if v := hc.resp.Header.Peek(unsetName1); len(v) > 0 {
		t.Errorf("httpClient.processHeaderRules() not unset header '%s'", unsetName1)
	}

	if v := hc.resp.Header.Peek(unsetName2); len(v) > 0 {
		t.Errorf("httpClient.processHeaderRules() not unset header '%s'", unsetName2)
	}
}
