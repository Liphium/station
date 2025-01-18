package studio_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/studio"
	"github.com/pion/webrtc/v4"
)

// Action: st_join
func joinStudio(c *pipeshandler.Context, offer webrtc.SessionDescription) pipes.Event {

	// Create a new client connection for the studio
	s := studio.GetStudio(c.Client.Session)
	answer, err := s.NewClientConnection(c.Client.ID, offer)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"answer":  answer,
	})
}
