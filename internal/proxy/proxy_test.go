package proxy

import (
	"fmt"
	"testing"
	"time"

	"github.com/savsgio/kratgo/internal/proxy/config"
	"github.com/valyala/fasthttp"
)

func testConfig() config.Config {
	return config.Config{
		Proxy: config.Proxy{
			Addr:        "localhost:8000",
			BackendAddr: "localhost:9997",
		},
		Cache: config.Cache{
			TTL: 30 * time.Second,
		},
		Invalidator: config.Invalidator{
			Addr:       "localhost:8001",
			MaxWorkers: 1,
		},
		LogLevel:  "fatal",
		LogOutput: "console",
	}
}

type mockHTTPClient struct{}

func (mock *mockHTTPClient) Do(req *fasthttp.Request, resp *fasthttp.Response) error {
	resp.SetBodyString("Kratgo Http Cache")
	resp.Header.Set("X-Data", "Bench")
	resp.SetStatusCode(200)

	return nil
}

func BenchmarkHandler(b *testing.B) {
	p, err := New(testConfig())
	if err != nil {
		panic(err)
	}
	p.httpClient = new(mockHTTPClient)

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.SetRequestURI("/bench")
	ctx.Request.Header.SetMethod("GET")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.handler(ctx)
	}
}

func BenchmarkHandlerWithoutCache(b *testing.B) {
	path := "/bench"
	cfg := testConfig()
	cfg.Proxy.Nocache = []string{
		fmt.Sprintf("$(path) == '%s'", path),
	}

	p, err := New(cfg)
	if err != nil {
		panic(err)
	}
	p.httpClient = new(mockHTTPClient)

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.SetRequestURI(path)
	ctx.Request.Header.SetMethod("GET")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.handler(ctx)
	}
}
