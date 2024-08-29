package space

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// Action: spc_leave
func leaveCall(ctx *pipeshandler.Context, data interface{}) pipes.Event {

	// Check if in space
	if !caching.IsInSpace(ctx.Client.ID) {
		return pipeshandler.ErrorResponse(ctx, "not.in.space", nil)
	}

	// Leave space
	valid := caching.LeaveSpace(ctx.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(ctx, localization.ErrorServer, nil)
	}

	// Send success
	return pipeshandler.SuccessResponse(ctx)
}
