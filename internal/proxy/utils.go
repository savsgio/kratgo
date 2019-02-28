package proxy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/savsgio/kratgo/internal/config"

	"github.com/savsgio/gotils"
	"github.com/valyala/fasthttp"
)

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

// TOOLS

func intSliceIndexOf(vs []int, t int) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func intSliceInclude(vs []int, t int) bool {
	return intSliceIndexOf(vs, t) >= 0
}

func stringSliceIndexOf(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func stringSliceInclude(vs []string, t string) bool {
	return stringSliceIndexOf(vs, t) >= 0
}

// HTTP

func cloneHeaders(dst, src *fasthttp.RequestHeader) {
	src.VisitAll(func(key, value []byte) {
		if !stringSliceInclude(hopHeaders, gotils.B2S(key)) {
			// fmt.Println(gotils.B2S(key), gotils.B2S(value))
			dst.SetCanonical(key, value)
		}
	})
}

func getEvalValue(req *fasthttp.Request, resp *fasthttp.Response, name, key string) string {
	value := name

	switch name {
	case config.EvalMethodVar:
		value = gotils.B2S(req.Header.Method())

	case config.EvalHostVar:
		value = gotils.B2S(req.Host())

	case config.EvalPathVar:
		value = gotils.B2S(req.URI().PathOriginal())

	case config.EvalContentTypeVar:
		value = gotils.B2S(resp.Header.ContentType())

	case config.EvalStatusCodeVar:
		value = strconv.Itoa(resp.StatusCode())

	default:
		if strings.HasPrefix(name, config.EvalReqHeaderVar) {
			value = gotils.B2S(req.Header.Peek(key))

		} else if strings.HasPrefix(name, config.EvalRespHeaderVar) {
			value = gotils.B2S(resp.Header.Peek(key))

		} else if strings.HasPrefix(name, config.EvalCookieVar) {
			value = gotils.B2S(req.Header.Cookie(key))
		}
	}

	return value
}

func checkIfNoCache(req *fasthttp.Request, resp *fasthttp.Response, rules []rule, params *evalParams) (bool, error) {
	for _, r := range rules {
		params.reset()

		for _, p := range r.params {
			params.set(p.name, getEvalValue(req, resp, p.name, p.subKey))
		}

		result, err := r.expr.Evaluate(params.all())
		if err != nil {
			return false, fmt.Errorf("Invalid nocache rule: %v", err)
		}

		if result.(bool) {
			return true, nil
		}
	}

	return false, nil
}

func processHeaderRules(req *fasthttp.Request, resp *fasthttp.Response, rules []headerRule, params *evalParams) error {
	for _, r := range rules {
		params.reset()

		executeHeaderRule := true

		if r.expr != nil {
			for _, p := range r.params {
				params.set(p.name, getEvalValue(req, resp, p.name, p.subKey))
			}

			result, err := r.expr.Evaluate(params.all())
			if err != nil {
				return err
			}

			executeHeaderRule = result.(bool)
		}

		if !executeHeaderRule {
			continue
		}

		if r.action == setHeaderAction {
			resp.Header.Set(r.name, getEvalValue(req, resp, r.value.value, r.value.subKey))
		} else {
			resp.Header.Del(r.name)
		}
	}

	return nil
}
