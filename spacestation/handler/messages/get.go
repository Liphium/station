package message_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

// Action: msg_get
func getMessage(c *pipeshandler.Context, id string) pipes.Event {

	// Get the messages with the specified id
	msg, err := caching.GetMessageById(c.Client.Session, id)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"message": msg,
	})
}
