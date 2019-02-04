package invalidator

import (
	"encoding/json"
	"fmt"

	"github.com/savsgio/gotils"
	"github.com/valyala/fasthttp"
)

func (i *Invalidator) httpHandler(ctx *fasthttp.RequestCtx) {
	path := gotils.B2S(ctx.Path())

	if i.log.DebugEnabled() {
		i.log.Debugf("%s - %s", ctx.Method(), path)
	}

	if !ctx.IsPost() {
		err := fmt.Errorf("Invalid request method: %s", ctx.Method())
		i.log.Info(err)
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	switch path {
	case "/invalidate/":
		entry := AcquireEntry()
		body := ctx.PostBody()

		if i.log.DebugEnabled() {
			i.log.Debugf("Invalidation received: %s", body)
		}

		err := json.Unmarshal(body, entry)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
			ReleaseEntry(entry)
			return
		}

		i.Add(*entry)

		ReleaseEntry(entry)

		ctx.SetBodyString("OK")

	default:
		ctx.NotFound()
	}

}
