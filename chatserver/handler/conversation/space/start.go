package space

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// Action: spc_start
func start(c *pipeshandler.Context, data interface{}) pipes.Event {

	// Create space
	roomId, appToken, err := caching.CreateSpace(c.Client.ID)
	if err != nil {
		util.Log.Println(err.Error())
		return pipeshandler.ErrorResponse(c, err.Error(), err)
	}

	// Send space info
	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"id":      roomId,
		"token":   appToken,
	})
}
