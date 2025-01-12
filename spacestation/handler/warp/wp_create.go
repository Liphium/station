package warp_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

// Action: wp_create
func create(c *pipeshandler.Context, port uint) pipes.Event {

	// Make sure the port is valid
	if port < 1024 || port > 65535 {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequestContent, nil)
	}

	// Create a new Warp for this port
	warpId, err := caching.NewWarp(c.Client.Session, c.Client.ID, port)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	// Send the warp id to the hoster so they can accept users of the Warp
	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"id":      warpId,
	})
}
