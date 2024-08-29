package space

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// Action: spc_start
func start(c *pipeshandler.Context, data interface{}) pipes.Event {

	/*
		TODO: Re-enable
		if caching.IsInSpace(message.Client.ID) {
			wshandler.ErrorResponse(message, "already.in.space")
			return
		}
	*/

	// Create space
	roomId, appToken, valid := caching.CreateSpace(c.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	// Send space info
	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"id":      roomId,
		"token":   appToken,
	})
}
