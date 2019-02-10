package proxy

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
	"github.com/savsgio/kratgo/internal/config"
	"github.com/valyala/fasthttp"
)

func Test_intSliceIndexOf(t *testing.T) {
	array := []int{1, 2, 3, 4, 5}

	n := 3
	if i := intSliceIndexOf(array, n); i < 0 {
		t.Errorf("intSliceIndexOf() = %v, want %v", i, 2)
	}

	n = 9
	if i := intSliceIndexOf(array, n); i > -1 {
		t.Errorf("intSliceIndexOf() = %v, want %v", i, -1)
	}
}

func Test_intSliceInclude(t *testing.T) {
	array := []int{1, 2, 3, 4, 5}

	n := 3
	if ok := intSliceInclude(array, n); !ok {
		t.Errorf("intSliceIndexOf() = %v, want %v", ok, true)
	}

	n = 9
	if ok := intSliceInclude(array, n); ok {
		t.Errorf("intSliceIndexOf() = %v, want %v", ok, false)
	}
}

func Test_stringSliceIndexOf(t *testing.T) {
	array := []string{"kratgo", "fast", "http", "cache"}

	s := "fast"
	if i := stringSliceIndexOf(array, s); i < 0 {
		t.Errorf("stringSliceIndexOf() = %v, want %v", i, 2)
	}

	s = "slow"
	if i := stringSliceIndexOf(array, s); i > -1 {
		t.Errorf("stringSliceIndexOf() = %v, want %v", i, -1)
	}
}

func Test_stringSliceInclude(t *testing.T) {
	array := []string{"kratgo", "fast", "http", "cache"}

	s := "fast"
	if ok := stringSliceInclude(array, s); !ok {
		t.Errorf("stringSliceInclude() = %v, want %v", ok, true)
	}

	s = "slow"
	if ok := stringSliceInclude(array, s); ok {
		t.Errorf("stringSliceInclude() = %v, want %v", ok, false)
	}
}

func Test_cloneHeaders(t *testing.T) {
	k1 := "Kratgo"
	v1 := "Fast"

	req1 := fasthttp.AcquireRequest()
	req2 := fasthttp.AcquireRequest()

	req1.Header.Set(k1, v1)
	for i, header := range hopHeaders {
		req1.Header.Set(header, strconv.Itoa(i))
	}

	cloneHeaders(&req2.Header, &req1.Header)

	isK1InReq2 := false
	req2.Header.VisitAll(func(k, v []byte) {
		if stringSliceInclude(hopHeaders, string(k)) {
			t.Errorf("cloneHeaders() invalid header '%s'", k)
		}

		if string(k) == k1 {
			isK1InReq2 = true

			if string(v) != v1 {
				t.Errorf("cloneHeaders() invalid header value of '%s' = '%s', want '%s'", k, v, v1)
			}
		}
	})

	if !isK1InReq2 {
		t.Errorf("cloneHeaders() the header '%s' is not cloned", k1)
	}
}

func Test_decodeResponseBody(t *testing.T) {
	type args struct {
		encodeType string
		body       string
	}

	body := "Kratgo always is fast"

	var b bytes.Buffer
	w1 := gzip.NewWriter(&b)
	w1.Write([]byte(body))
	w1.Close()
	bodyGZIp := b.String()

	b.Reset()

	w2 := zlib.NewWriter(&b)
	w2.Write([]byte(body))
	w2.Close()
	bodyDeflate := b.String()

	tests := []struct {
		name string
		args args
	}{
		{
			name: "gzip",
			args: args{
				encodeType: "gzip",
				body:       bodyGZIp,
			},
		},
		{
			name: "deflate",
			args: args{
				encodeType: "deflate",
				body:       bodyDeflate,
			},
		},
		{
			name: "raw",
			args: args{
				encodeType: "raw",
				body:       body,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := new(fasthttp.RequestCtx)
			resp := fasthttp.AcquireResponse()
			resp.SetBodyString(tt.args.body)
			resp.Header.Set(headerContentEncoding, tt.args.encodeType)

			err := decodeResponseBody(resp, ctx)
			if err != nil {
				t.Fatal(err)
			}

			respBody := string(ctx.Response.Body())
			if respBody != body {
				t.Errorf("decodeResponseBody() body = '%s', want '%s'", respBody, body)
			}

		})
	}
}

func Test_getEvalValue(t *testing.T) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	method := "POST"
	host := "www.kratgo.com"
	path := "/data/"
	contentType := "application/json"
	statusCode := 301
	reqHeaderName := "X-Kratgo"
	reqHeaderValue := "Fast"
	respHeaderName := "X-Data"
	respHeaderValue := "false"
	cookieName := "kratcookie"
	cookieValue := "1234"

	req.Header.SetMethod(method)
	req.Header.SetHost(host)
	req.SetRequestURI(path)
	resp.Header.SetContentType(contentType)
	resp.SetStatusCode(statusCode)
	req.Header.Set(reqHeaderName, reqHeaderValue)
	resp.Header.Set(respHeaderName, respHeaderValue)
	req.Header.SetCookie(cookieName, cookieValue)

	type args struct {
		name string
		key  string
	}

	type want struct {
		value string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "method",
			args: args{
				name: config.EvalMethodVar,
			},
			want: want{
				value: method,
			},
		},
		{
			name: "host",
			args: args{
				name: config.EvalHostVar,
			},
			want: want{
				value: host,
			},
		},
		{
			name: "path",
			args: args{
				name: config.EvalPathVar,
			},
			want: want{
				value: path,
			},
		},
		{
			name: "content-type",
			args: args{
				name: config.EvalContentTypeVar,
			},
			want: want{
				value: contentType,
			},
		},
		{
			name: "status-code",
			args: args{
				name: config.EvalStatusCodeVar,
			},
			want: want{
				value: strconv.Itoa(statusCode),
			},
		},
		{
			name: "request-header",
			args: args{
				name: config.EvalReqHeaderVar,
				key:  reqHeaderName,
			},
			want: want{
				value: reqHeaderValue,
			},
		},
		{
			name: "response-header",
			args: args{
				name: config.EvalRespHeaderVar,
				key:  respHeaderName,
			},
			want: want{
				value: respHeaderValue,
			},
		},
		{
			name: "cookie",
			args: args{
				name: config.EvalCookieVar,
				key:  cookieName,
			},
			want: want{
				value: cookieValue,
			},
		},
		{
			name: "unknown",
			args: args{
				name: "unknown",
			},
			want: want{
				value: "unknown",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEvalValue(req, resp, tt.args.name, tt.args.key); got != tt.want.value {
				t.Errorf("getEvalValue() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func Test_checkIfNoCache(t *testing.T) {
	cfg := testConfig()
	cfg.FileConfig.Nocache = []string{
		"$(method) == 'POST' && $(host) != 'www.kratgo.com'",
	}
	p, _ := New(cfg)
	p.parseNocacheRules()

	type args struct {
		method string
		host   string
	}

	type want struct {
		noCache bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Yes",
			args: args{
				method: "POST",
				host:   "www.example.com",
			},
			want: want{
				noCache: true,
			},
		},
		{
			name: "No1",
			args: args{
				method: "GET",
				host:   "www.kratgo.com",
			},
			want: want{
				noCache: false,
			},
		},
		{
			name: "No2",
			args: args{
				method: "POST",
				host:   "www.kratgo.com",
			},
			want: want{
				noCache: false,
			},
		},
		{
			name: "No3",
			args: args{
				method: "GET",
				host:   "www.example.com",
			},
			want: want{
				noCache: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()
			params := acquireEvalParams()

			req.Header.SetMethod(tt.args.method)
			req.Header.SetHost(tt.args.host)

			noCache, err := checkIfNoCache(req, resp, p.nocacheRules, params)
			if err != nil {
				t.Fatal(err)
			}

			if noCache != tt.want.noCache {
				t.Errorf("checkIfNoCache() = '%v', want '%v'", noCache, tt.want.noCache)
			}
		})
	}
}
