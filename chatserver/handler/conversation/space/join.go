package space

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipeshandler"
)

// Action: spc_join
func joinCall(ctx pipeshandler.Context) {

	if ctx.ValidateForm("id") {
		pipeshandler.ErrorResponse(ctx, localization.InvalidRequest)
		return
	}

	if caching.IsInSpace(ctx.Client.ID) {
		pipeshandler.ErrorResponse(ctx, "already.in.space")
		return
	}

	// Create space
	appToken, valid := caching.JoinSpace(ctx.Client.ID, ctx.Data["id"].(string))
	if !valid {
		pipeshandler.ErrorResponse(ctx, localization.ErrorServer)
		return
	}

	// Send space info
	pipeshandler.NormalResponse(ctx, map[string]interface{}{
		"success": true,
		"token":   appToken,
	})
}
