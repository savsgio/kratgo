package proxy

import (
	"sync"

	"github.com/valyala/fasthttp"
)

var httpClientPool = sync.Pool{
	New: func() interface{} {
		return &httpClient{
			req:  fasthttp.AcquireRequest(),
			resp: fasthttp.AcquireResponse(),
		}
	},
}

func acquireHTTPClient() *httpClient {
	return httpClientPool.Get().(*httpClient)
}

func releaseHTTPClient(hc *httpClient) {
	hc.reset()
	httpClientPool.Put(hc)
}

func (hc *httpClient) reset() {
	hc.executeHeaderRule = false

	hc.req.Reset()
	hc.resp.Reset()
}

func (hc *httpClient) do(f fetcher) error {
	hc.req.Header.Set(clientReqHeaderKey, clientReqHeaderValue)

	return f.Do(hc.req, hc.resp)
}

func (hc *httpClient) setMethodBytes(method []byte) {
	hc.req.Header.SetMethodBytes(method)
}

func (hc *httpClient) setRequestURIBytes(uri []byte) {
	hc.req.SetRequestURIBytes(uri)
}

func (hc *httpClient) setRequestBody(body []byte) {
	hc.req.SetBody(body)
}

func (hc *httpClient) copyReqHeaderTo(h *fasthttp.RequestHeader) {
	hc.req.Header.CopyTo(h)
}

func (hc *httpClient) copyRespHeaderTo(h *fasthttp.ResponseHeader) {
	hc.resp.Header.CopyTo(h)
}

func (hc *httpClient) respHeaderPeek(key string) []byte {
	return hc.resp.Header.Peek(key)
}

func (hc *httpClient) statusCode() int {
	return hc.resp.StatusCode()
}

func (hc *httpClient) body() []byte {
	return hc.resp.Body()
}

func (hc *httpClient) processHeaderRules(rules []headerRule, params *evalParams) error {
	for _, r := range rules {
		params.reset()

		hc.executeHeaderRule = true

		if r.expr != nil {
			for _, p := range r.params {
				params.set(p.name, getEvalValue(hc.req, hc.resp, p.name, p.subKey))
			}

			result, err := r.expr.Evaluate(params.all())
			if err != nil {
				return err
			}

			hc.executeHeaderRule = result.(bool)
		}

		if !hc.executeHeaderRule {
			continue
		}

		if r.action == setHeaderAction {
			hc.resp.Header.Set(r.name, getEvalValue(hc.req, hc.resp, r.value.value, r.value.subKey))
		} else {
			hc.resp.Header.Del(r.name)
		}
	}

	return nil
}
