package studio_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/studio"
)

// Action: st_leave
func leaveStudio(c *pipeshandler.Context, _ interface{}) pipes.Event {

	// Disconnect the client
	s := studio.GetStudio(c.Client.Session)
	if err := s.Disconnect(c.Client.ID); err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	return pipeshandler.SuccessResponse(c)
}
