package townsquare_handlers

import (
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// Action: townsquare_send
func sendMessage(ctx pipeshandler.Context) {

	if ctx.ValidateForm("content", "attachments") {
		pipeshandler.ErrorResponse(ctx, localization.InvalidRequest)
		return
	}

	id := caching.TownsquareMessageId()

	message := caching.TownsquareMessage{
		ID:          id,
		Sender:      ctx.Client.ID,
		Attachments: ctx.Data["attachments"].(string),
		Content:     ctx.Data["content"].(string),
		Timestamp:   time.Now().UnixMilli(),
	}

	caching.SendTownsquareMessageEvent(pipes.Event{
		Name: "ts_message",
		Data: map[string]interface{}{
			"msg": message,
		},
	})

	pipeshandler.SuccessResponse(ctx)
}

// Action: townsquare_delete
func deleteMessage(context pipeshandler.Context) {

}
