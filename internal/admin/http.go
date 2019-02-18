package admin

import (
	"github.com/savsgio/kratgo/internal/invalidator"

	json "github.com/json-iterator/go"
	"github.com/savsgio/atreugo/v7"
)

func (a *Admin) invalidateView(ctx *atreugo.RequestCtx) error {
	entry := invalidator.AcquireEntry()
	body := ctx.PostBody()

	if a.log.DebugEnabled() {
		a.log.Debugf("Invalidation received: %s", body)
	}

	err := json.Unmarshal(body, entry)
	if err != nil {
		invalidator.ReleaseEntry(entry)
		return err
	}

	if err = a.invalidator.Add(*entry); err != nil {
		a.log.Errorf("Could not add a invalidation entry '%s': %v", body, err)
		invalidator.ReleaseEntry(entry)
		return ctx.TextResponse(err.Error(), 400)
	}

	invalidator.ReleaseEntry(entry)

	return ctx.TextResponse("OK")
}
