package space

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipeshandler/wshandler"
)

// Action: spc_leave
func leaveCall(message wshandler.Message) {

	// Check if in space
	if !caching.IsInSpace(message.Client.ID) {
		wshandler.ErrorResponse(message, "not.in.space")
		return
	}

	// Leave space
	valid := caching.LeaveSpace(message.Client.ID)
	if !valid {
		wshandler.ErrorResponse(message, localization.ErrorServer)
		return
	}

	// Send success
	wshandler.SuccessResponse(message)
}
