package studio_track_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/studio"
)

// Action: st_tr_subscribe
func subscribeToTrack(c *pipeshandler.Context, action struct {
	Track   string `json:"track"`
	Channel string `json:"channel"`
}) pipes.Event {
	s := studio.GetStudio(c.Client.Session)

	// Get the client
	client, valid := s.GetClient(c.Client.ID)
	if valid {
		return pipeshandler.ErrorResponse(c, localization.ErrorDidntJoinStudio, nil)
	}

	// Get the track
	track, valid := s.GetTrack(action.Track)
	if !valid {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequestContent, nil)
	}

	// Subscribe to the track
	if err := track.NewSubscription(client, action.Channel); err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	return pipeshandler.SuccessResponse(c)
}
