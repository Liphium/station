package space

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipeshandler/wshandler"
)

// Action: spc_join
func joinCall(message wshandler.Message) {

	if message.ValidateForm("id") {
		wshandler.ErrorResponse(message, localization.InvalidRequest)
		return
	}

	if caching.IsInSpace(message.Client.ID) {
		wshandler.ErrorResponse(message, "already.in.space")
		return
	}

	// Create space
	appToken, valid := caching.JoinSpace(message.Client.ID, message.Data["id"].(string), integration.ClusterID)
	if !valid {
		wshandler.ErrorResponse(message, localization.ErrorServer)
		return
	}

	// Send space info
	wshandler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"token":   appToken,
	})
}
