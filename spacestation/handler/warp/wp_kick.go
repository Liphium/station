package warp_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

// Action: wp_kick
func kick(c *pipeshandler.Context, action struct {
	Warp   string `json:"w"`
	Target string `json:"t"`
}) pipes.Event {

	// Get the Warp related to the packet
	warp, err := caching.GetWarp(c.Client.Session, action.Warp)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequestContent, err)
	}

	// Make sure it's the hoster sending the event
	if warp.Hoster != c.Client.ID {
		return pipeshandler.ErrorResponse(c, localization.ErrorNoPermission, err)
	}

	// Remove from the list of Warp receivers
	if err := caching.RemoveClientFromWarp(action.Target, c.Client.Session, action.Warp); err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequestContent, err)
	}

	return pipeshandler.SuccessResponse(c)
}
