package proxy

import (
	"fmt"
	"regexp"
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

func TestProxy_newEvaluableExpression(t *testing.T) {
	type args struct {
		rule string
	}

	type want struct {
		strExpr     string
		regexExpr   *regexp.Regexp
		totalParams int
		err         bool
	}

	p, err := New(testConfig())
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "version",
			args: args{
				rule: fmt.Sprintf("$(version) == '%s'", version),
			},
			want: want{
				strExpr:     fmt.Sprintf("%s == '%s'", config.EvalVersionVar, version),
				totalParams: 1,
				err:         false,
			},
		},
		{
			name: "method",
			args: args{
				rule: fmt.Sprintf("$(method) == '%s'", "GET"),
			},
			want: want{
				strExpr:     fmt.Sprintf("%s == '%s'", config.EvalMethodVar, "GET"),
				totalParams: 1,
				err:         false,
			},
		},
		{
			name: "host",
			args: args{
				rule: fmt.Sprintf("$(host) == '%s'", "www.kratgo.com"),
			},
			want: want{
				strExpr:     fmt.Sprintf("%s == '%s'", config.EvalHostVar, "www.kratgo.com"),
				totalParams: 1,
				err:         false,
			},
		},
		{
			name: "path",
			args: args{
				rule: fmt.Sprintf("$(path) == '%s'", "/es/"),
			},
			want: want{
				strExpr:     fmt.Sprintf("%s == '%s'", config.EvalPathVar, "/es/"),
				totalParams: 1,
				err:         false,
			},
		},
		{
			name: "contentType",
			args: args{
				rule: fmt.Sprintf("$(contentType) == '%s'", "text/html"),
			},
			want: want{
				strExpr:     fmt.Sprintf("%s == '%s'", config.EvalContentTypeVar, "text/html"),
				totalParams: 1,
				err:         false,
			},
		},
		{
			name: "statusCode",
			args: args{
				rule: fmt.Sprintf("$(statusCode) == '%s'", "200"),
			},
			want: want{
				strExpr:     fmt.Sprintf("%s == '%s'", config.EvalStatusCodeVar, "200"),
				totalParams: 1,
				err:         false,
			},
		},
		{
			name: "req.header::<NAME>",
			args: args{
				rule: fmt.Sprintf("$(req.header::X-Data) == '%s'", "Kratgo"),
			},
			want: want{
				regexExpr:   regexp.MustCompile(fmt.Sprintf("%s([0-9]{2}) == '%s'", config.EvalReqHeaderVar, "Kratgo")),
				totalParams: 1,
				err:         false,
			},
		},
		{
			name: "resp.header::<NAME>",
			args: args{
				rule: fmt.Sprintf("$(resp.header::X-Data) == '%s'", "Kratgo"),
			},
			want: want{
				regexExpr:   regexp.MustCompile(fmt.Sprintf("%s([0-9]{2}) == '%s'", config.EvalRespHeaderVar, "Kratgo")),
				totalParams: 1,
				err:         false,
			},
		},
		{
			name: "cookie::<NAME>",
			args: args{
				rule: fmt.Sprintf("$(cookie::X-Data) == '%s'", "Kratgo"),
			},
			want: want{
				regexExpr:   regexp.MustCompile(fmt.Sprintf("%s([0-9]{2}) == '%s'", config.EvalCookieVar, "Kratgo")),
				totalParams: 1,
				err:         false,
			},
		},
		{
			name: "combo",
			args: args{
				rule: fmt.Sprintf("$(path) == '%s' && $(method) != '%s'", "/kratgo", "GET"),
			},
			want: want{
				strExpr:     fmt.Sprintf("%s == '%s' && %s != '%s'", config.EvalPathVar, "/kratgo", config.EvalMethodVar, "GET"),
				totalParams: 2,
				err:         false,
			},
		},
		{
			name: "Error",
			args: args{
				rule: "$(test) /() thod) != asdasd3'",
			},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, params, err := p.newEvaluableExpression(tt.args.rule)

			if (tt.want.err && err == nil) || (!tt.want.err && err != nil) {
				t.Fatalf("Proxy.newEvaluableExpression() returns error '%v', want error '%v'", err, tt.want.err)
			}

			if !tt.want.err {
				strExpr := expr.String()
				if tt.want.regexExpr != nil {
					if !tt.want.regexExpr.MatchString(strExpr) {
						t.Errorf("Proxy.newEvaluableExpression() = '%s', want '%s'", strExpr, tt.want.regexExpr.String())
					}
				} else {
					if strExpr != tt.want.strExpr {
						t.Errorf("Proxy.newEvaluableExpression() = '%s', want '%s'", expr.String(), tt.want.strExpr)
					}
				}

				if len(params) != tt.want.totalParams {
					t.Errorf("Proxy.newEvaluableExpression() len(params) = '%d', want '%d'", len(params), tt.want.totalParams)
				}
			}

		})
	}
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
