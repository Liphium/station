package space

import (
	"os"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// Action: spc_start
func start(c *pipeshandler.Context, data interface{}) pipes.Event {

	if os.Getenv("SPACES_APP") == "" {
		return pipeshandler.ErrorResponse(c, localization.ErrorSpacesNotSetup, nil)
	}

	// Create space
	roomId, appToken, err := caching.CreateSpace(c.Client.ID)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorSpacesNotSetup, err)
	}

	// Send space info
	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"id":      roomId,
		"token":   appToken,
	})
}
