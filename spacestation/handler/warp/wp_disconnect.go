package warp_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

// Action: wp_disconnect
func disconnect(c *pipeshandler.Context, warpId string) pipes.Event {

	// Remove from the list of Warp receivers
	if err := caching.RemoveClientFromWarp(c.Client.ID, c.Client.Session, warpId); err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequestContent, err)
	}

	return pipeshandler.SuccessResponse(c)
}
