package message_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

// Action: msg_list_after
func listMessagesAfter(c *pipeshandler.Context, time int64) pipes.Event {

	// Get all the messages after the specified time
	messages, err := caching.GetMessagesAfter(c.Client.Session, time)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success":  true,
		"messages": messages,
	})
}
