package townsquare_handlers

import (
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

type sendMessageAction struct {
	Content     string `json:"content"`
	Attachments string `json:"attachments"`
}

// Action: townsquare_send
func sendMessage(c *pipeshandler.Context, action sendMessageAction) pipes.Event {

	id := caching.TownsquareMessageId()

	message := caching.TownsquareMessage{
		ID:          id,
		Sender:      c.Client.ID,
		Attachments: action.Attachments,
		Content:     action.Content,
		Timestamp:   time.Now().UnixMilli(),
	}

	caching.SendTownsquareMessageEvent(pipes.Event{
		Name: "ts_message",
		Data: map[string]interface{}{
			"msg": message,
		},
	})

	return pipeshandler.SuccessResponse(c)
}

// Action: townsquare_delete
func deleteMessage(c *pipeshandler.Context, data interface{}) pipes.Event {
	return pipeshandler.SuccessResponse(c)
}
