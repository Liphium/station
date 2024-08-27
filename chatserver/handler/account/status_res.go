package account

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util/localization"
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

	// Make sure it's a private conversation
	if len(members) > 2 {
		pipeshandler.ErrorResponse(ctx, localization.InvalidRequest)
		return
	}

	// Get the other member to send the status to
	var otherMember caching.StoredMember
	for _, member := range members {
		if member.TokenID != convToken.ID {
			otherMember = member
		}
	}

	// Send the event
	if err := caching.SendEventToMembers([]caching.StoredMember{otherMember}, StatusEvent(status, data, convToken.Conversation, convToken.ID, ":a")); err != nil {
		pipeshandler.ErrorResponse(ctx, localization.ErrorServer)
		return
	}

	pipeshandler.SuccessResponse(ctx)
}
