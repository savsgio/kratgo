package admin

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/savsgio/gotils"
	"github.com/savsgio/kratgo/internal/invalidator"
	"github.com/valyala/fasthttp"
)

func (a *Admin) httpHandler(ctx *fasthttp.RequestCtx) {
	path := gotils.B2S(ctx.Path())

	if a.log.DebugEnabled() {
		a.log.Debugf("%s - %s", ctx.Method(), path)
	}

	if !ctx.IsPost() {
		err := fmt.Errorf("Invalid request method: %s", ctx.Method())
		a.log.Info(err)
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	switch path {
	case "/invalidate/":
		entry := invalidator.AcquireEntry()
		body := ctx.PostBody()

		if a.log.DebugEnabled() {
			a.log.Debugf("Invalidation received: %s", body)
		}

		err := jsoniter.Unmarshal(body, entry)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
			invalidator.ReleaseEntry(entry)
			return
		}

		a.invalidator.Add(*entry)

		invalidator.ReleaseEntry(entry)

		ctx.SetBodyString("OK")

	default:
		ctx.NotFound()
	}

}
