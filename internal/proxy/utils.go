package proxy

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/savsgio/gotils"
	"github.com/savsgio/kratgo/internal/proxy/config"
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

// LOG

func getLogOutput(output string) (*os.File, error) {
	if output != "" && output != "console" {
		dirPath, _ := path.Split(output)
		if err := os.MkdirAll(dirPath, os.ModeDir); err != nil {
			return nil, err
		}
		return os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0755)
	}

	return os.Stderr, nil
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

func decodeResponseBody(resp *fasthttp.Response, ctx *fasthttp.RequestCtx) error {
	encoding := gotils.B2S(resp.Header.Peek(headerContentEncoding))
	switch encoding {
	case "gzip":
		body, err := resp.BodyGunzip()
		if err != nil {
			return err
		}
		ctx.SetBody(body)

	case "deflate":
		body, err := resp.BodyInflate()
		if err != nil {
			return err
		}
		ctx.SetBody(body)

	default:
		resp.BodyWriteTo(ctx)
	}

	return nil
}

func getEvalValue(req *fasthttp.Request, resp *fasthttp.Response, name, key string) string {
	value := name

	switch name {
	case config.EvalVersionVar:
		value = version

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

func checkIfNoCache(req *fasthttp.Request, resp *fasthttp.Response, rules []Rule, params *evalParams) (bool, error) {
	for _, r := range rules {
		params.reset()

		for _, p := range r.params {
			params.set(p.name, getEvalValue(req, resp, p.name, p.subKey))
		}

		result, err := r.expr.Evaluate(params.p)
		if err != nil {
			return false, fmt.Errorf("Invalid nocache rule: %v", err)
		}

		if result.(bool) {
			return true, nil
		}
	}

	return false, nil
}
