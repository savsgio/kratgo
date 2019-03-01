package proxy

import (
	"strconv"
	"testing"

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

func Test_getEvalValue(t *testing.T) {
	ctx := new(fasthttp.RequestCtx)

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

	ctx.Request.SetRequestURI(path)
	ctx.Request.Header.SetMethod(method)
	ctx.Request.Header.SetHost(host)
	ctx.Request.Header.Set(reqHeaderName, reqHeaderValue)
	ctx.Request.Header.SetCookie(cookieName, cookieValue)

	ctx.Response.Header.SetContentType(contentType)
	ctx.Response.Header.Set(respHeaderName, respHeaderValue)
	ctx.Response.SetStatusCode(statusCode)

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
			if got := getEvalValue(ctx, tt.args.name, tt.args.key); got != tt.want.value {
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

	type args struct {
		method        string
		host          string
		delRuleParams bool
	}

	type want struct {
		noCache bool
		err     bool
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
				err:     false,
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
				err:     false,
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
				err:     false,
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
				err:     false,
			},
		},
		{
			name: "Error",
			args: args{
				method:        "GET",
				host:          "www.example.com",
				delRuleParams: true,
			},
			want: want{
				noCache: false,
				err:     true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := new(fasthttp.RequestCtx)
			params := acquireEvalParams()

			ctx.Request.Header.SetMethod(tt.args.method)
			ctx.Request.Header.SetHost(tt.args.host)

			if tt.args.delRuleParams {
				for i := range p.nocacheRules {
					r := &p.nocacheRules[i]
					r.params = r.params[:0]
				}
			}

			noCache, err := checkIfNoCache(ctx, p.nocacheRules, params)
			if (err != nil) != tt.want.err {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.want.err {
				return
			}

			if noCache != tt.want.noCache {
				t.Errorf("checkIfNoCache() = '%v', want '%v'", noCache, tt.want.noCache)
			}
		})
	}
}

func TestHTTPClient_processHeaderRules(t *testing.T) {
	type args struct {
		processWithoutRuleParams bool
	}

	type want struct {
		err bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Ok",
			args: args{
				processWithoutRuleParams: false,
			},
			want: want{
				err: false,
			},
		},
		{
			name: "Error",
			args: args{
				processWithoutRuleParams: true,
			},
			want: want{
				err: true,
			},
		},
	}

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

	cfg := testConfig()
	cfg.FileConfig.Response.Headers.Set = setHeadersRulesConfig
	cfg.FileConfig.Response.Headers.Unset = unsetHeadersRulesConfig
	p, _ := New(cfg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.processWithoutRuleParams {
				for i := range p.headersRules {
					p.headersRules[i].params = p.headersRules[i].params[:0]
				}
			}

			params := acquireEvalParams()

			ctx := new(fasthttp.RequestCtx)

			ctx.Request.Header.Set("X-Data", "123")
			ctx.Response.Header.Set(unsetName1, "data")
			ctx.Response.Header.Set("FakeHeader", "fake data")
			ctx.Response.Header.Set("Content-Type", "text/html")

			err := processHeaderRules(ctx, p.headersRules, params)
			if (err != nil) != tt.want.err {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.want.err {
				return
			}

			if v := ctx.Response.Header.Peek(setName1); string(v) != setValue1 {
				t.Errorf("httpClient.processHeaderRules() not set header '%s' with value '%s', want '%s==%s'",
					setName1, setValue1, setName1, v)
			}

			if v := ctx.Response.Header.Peek(setName2); string(v) != setValue2 {
				t.Errorf("httpClient.processHeaderRules() not set header '%s' with value '%s', want '%s==%s'",
					setName2, setValue2, setName2, v)
			}

			if v := ctx.Response.Header.Peek(setName3); len(v) > 0 {
				t.Errorf("httpClient.processHeaderRules() header '%s' is setted but not fulfill the condition", setName3)
			}

			if v := ctx.Response.Header.Peek(unsetName1); len(v) > 0 {
				t.Errorf("httpClient.processHeaderRules() not unset header '%s'", unsetName1)
			}

			if v := ctx.Response.Header.Peek(unsetName2); len(v) > 0 {
				t.Errorf("httpClient.processHeaderRules() not unset header '%s'", unsetName2)
			}
		})
	}
}
