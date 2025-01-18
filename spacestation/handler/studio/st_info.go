package studio_handlers

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/studio"
)

// Action: st_info
func getStudioInfo(c *pipeshandler.Context, _ interface{}) pipes.Event {

	// Return all important info regarding studio
	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"stun":    studio.DefaultStunServer,
		"port":    studio.Port,
	})
}
