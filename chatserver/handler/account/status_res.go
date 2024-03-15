package account

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// Action: st_res
func respondToStatus(ctx pipeshandler.Context) {

	if ctx.ValidateForm("id", "token", "status", "data") {
		pipeshandler.ErrorResponse(ctx, localization.InvalidRequest)
		return
	}

	id := ctx.Data["id"].(string)
	token := ctx.Data["token"].(string)
	status := ctx.Data["status"].(string)
	data := ctx.Data["data"].(string)

	// Get from cache
	convToken, err := caching.ValidateToken(id, token)
	if err != nil {
		pipeshandler.ErrorResponse(ctx, localization.InvalidRequest)
		return
	}

	// Check if this is a valid conversation
	members, err := caching.LoadMembers(convToken.Conversation)
	if err != nil {
		pipeshandler.ErrorResponse(ctx, localization.ErrorServer)
		return
	}

	ids, nodes := caching.MembersToPipes(members)

	// Send the subscription event
	err = caching.Node.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.Conversation(ids, nodes),
		Event:   statusEvent(status, data, ":a"),
	})
	if err != nil {
		pipeshandler.ErrorResponse(ctx, localization.ErrorServer)
		return
	}

	pipeshandler.SuccessResponse(ctx)
}
