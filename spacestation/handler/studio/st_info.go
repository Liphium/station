package studio_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/studio"
)

// Action: st_info
func getStudioInfo(c *pipeshandler.Context, _ interface{}) pipes.Event {

	// Only return something in case Studio is enabled
	if !studio.Enabled {
		return pipeshandler.ErrorResponse(c, localization.ErrorStudioNotSupported, nil)
	}

	// Return all important info regarding studio
	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"stun":    studio.DefaultStunServer,
		"port":    studio.Port,
	})
}
