package admin

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/savsgio/atreugo/v7"
	"github.com/savsgio/kratgo/internal/invalidator"
)

func (a *Admin) invalidateView(ctx *atreugo.RequestCtx) error {
	entry := invalidator.AcquireEntry()
	body := ctx.PostBody()

	if a.log.DebugEnabled() {
		a.log.Debugf("Invalidation received: %s", body)
	}

	err := jsoniter.Unmarshal(body, entry)
	if err != nil {
		invalidator.ReleaseEntry(entry)
		return err
	}

	a.invalidator.Add(*entry)

	invalidator.ReleaseEntry(entry)

	return ctx.TextResponse("OK")
}
