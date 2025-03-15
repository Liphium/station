package studio_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/studio"
	"github.com/pion/webrtc/v4"
)

// Action: st_ice
func handleIceCandidate(c *pipeshandler.Context, candidate webrtc.ICECandidateInit) pipes.Event {

	// Create a new client connection for the studio
	s := studio.GetStudio(c.Client.Session)
	client, valid := s.GetClient(c.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequest, nil)
	}

	// Send the ice candidate to the client
	if err := client.HandleIceCandidate(candidate); err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	return pipeshandler.SuccessResponse(c)
}
