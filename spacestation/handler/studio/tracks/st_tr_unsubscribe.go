package studio_track_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/studio"
)

// Action: st_tr_unsubscribe
func unsubscribeToTrack(c *pipeshandler.Context, action struct {
	Track string `json:"track"`
}) pipes.Event {
	s := studio.GetStudio(c.Client.Session)

	// Get the client
	client, valid := s.GetClient(c.Client.ID)
	if valid {
		return pipeshandler.ErrorResponse(c, localization.ErrorDidntJoinStudio, nil)
	}

	// Get the subscription
	sub, valid := client.GetSubscription(action.Track)
	if !valid {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequestContent, nil)
	}
	sub.Delete()

	return pipeshandler.SuccessResponse(c)
}