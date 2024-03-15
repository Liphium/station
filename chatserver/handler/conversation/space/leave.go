package space

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipeshandler"
)

// Action: spc_leave
func leaveCall(ctx pipeshandler.Context) {

	// Check if in space
	if !caching.IsInSpace(ctx.Client.ID) {
		pipeshandler.ErrorResponse(ctx, "not.in.space")
		return
	}

	// Leave space
	valid := caching.LeaveSpace(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, localization.ErrorServer)
		return
	}

	// Send success
	pipeshandler.SuccessResponse(ctx)
}
