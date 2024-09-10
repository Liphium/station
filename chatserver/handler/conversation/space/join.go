package space

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// Action: spc_join
func joinCall(c *pipeshandler.Context, id string) pipes.Event {

	if caching.IsInSpace(c.Client.ID) {
		return pipeshandler.ErrorResponse(c, localization.ErrorAlreadyInSpace, nil)
	}

	// Create space
	appToken, valid := caching.JoinSpace(c.Client.ID, id)
	if !valid {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	// Send space info
	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"token":   appToken,
	})
}
