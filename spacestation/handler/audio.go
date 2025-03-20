package handler

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

// Action: set_muted
func setMuted(c *pipeshandler.Context, muted bool) pipes.Event {

	// Update the member data
	if !caching.UpdateMemberData(c.Client.Session, c.Client.ID, nil, caching.Ptr(muted), nil) {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	return pipeshandler.SuccessResponse(c)
}

// Action: set_deafened
func setDeafened(c *pipeshandler.Context, deafened bool) pipes.Event {

	// Update the member data
	if !caching.UpdateMemberData(c.Client.Session, c.Client.ID, nil, nil, caching.Ptr(deafened)) {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	return pipeshandler.SuccessResponse(c)
}
