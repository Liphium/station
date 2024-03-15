package space

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipeshandler"
)

// Action: spc_start
func start(ctx pipeshandler.Context) {

	/*
		TODO: Re-enable
		if caching.IsInSpace(message.Client.ID) {
			wshandler.ErrorResponse(message, "already.in.space")
			return
		}
	*/

	// Create space
	roomId, appToken, valid := caching.CreateSpace(ctx.Client.ID, integration.ClusterID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, localization.ErrorServer)
		return
	}

	// Send space info
	pipeshandler.NormalResponse(ctx, map[string]interface{}{
		"success": true,
		"id":      roomId,
		"token":   appToken,
	})
}
