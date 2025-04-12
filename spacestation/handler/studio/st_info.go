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
	body := map[string]interface{}{
		"success": true,
		"stun":    studio.StunServer,
		"port":    studio.Port,
	}
	if studio.TurnServer != "" {
		body["turn"] = studio.TurnServer
		body["turn_user"] = studio.TurnUsername
		body["turn_pass"] = studio.TurnPassword
	}
	return pipeshandler.NormalResponse(c, body)
}
