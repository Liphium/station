package handler

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

// Action: set_audio_state
func setAudioState(c *pipeshandler.Context, action struct {
	Muted    *bool `json:"muted"`
	Deafened *bool `json:"deafened"`
}) pipes.Event {

	// Make sure the user actually wants to change something
	if action.Muted == nil && action.Deafened == nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequest, nil)
	}

	// Update the member data
	if !caching.UpdateMemberData(c.Client.Session, c.Client.ID, nil, action.Muted, action.Deafened) {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	return pipeshandler.SuccessResponse(c)
}
