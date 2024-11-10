package message_handlers

import (
	"time"

	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

// Action: msg_get
func deleteMessage(c *pipeshandler.Context, id string) pipes.Event {

	// Delete the messages with the specified id
	err := caching.DeleteMessage(c.Client.Session, id)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	// Send an event to tell everyone about the deleted message
	contentJson, err := sonic.MarshalString(map[string]interface{}{
		"c": "msg.deleted",
		"a": []string{id},
	})
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}
	message := caching.Message{
		ID:           uuid.NewString(),
		Conversation: c.Client.Session + "@" + integration.Domain,
		Creation:     time.Now().UnixMilli(),
		Data:         contentJson,
		Edited:       false,
		Sender:       systemSender,
	}
	SendEventToMembers(c.Client.Session, pipes.Event{
		Name: "msg",
		Data: map[string]interface{}{
			"msg": message,
		},
	})

	return pipeshandler.SuccessResponse(c)
}
