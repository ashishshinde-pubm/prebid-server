package katana

import (
	"context"
	"time"

	"github.com/prebid/prebid-server/v2/hooks/hookstage"
	"github.com/prebid/prebid-server/v2/modules/prebid/katana/models"
)

func (k Katana) handleEntrypointHook(
	ctx context.Context,
	miCtx hookstage.ModuleInvocationContext,
	payload hookstage.EntrypointPayload,
) (result hookstage.HookResult[hookstage.EntrypointPayload], err error) {
	rctx := models.RequestCtx{}

	defer func() {
		result.ModuleContext = make(hookstage.ModuleContext)
		result.ModuleContext["rctx"] = rctx
	}()

	rctx = models.RequestCtx{
		StartTime: time.Now().Unix(),
		UA:        payload.Request.Header.Get("User-Agent"),
		ImpCtx:    make(map[string]models.ImpCtx),
	}

	return result, nil
}
