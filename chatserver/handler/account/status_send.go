package account

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/fetching"
	"github.com/Liphium/station/chatserver/handler/conversation"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// Action: st_send
func sendStatus(ctx pipeshandler.Context) {

	if ctx.ValidateForm("tokens", "status", "data") {
		pipeshandler.ErrorResponse(ctx, localization.InvalidRequest)
		return
	}

	// Save in database
	statusMessage := ctx.Data["status"].(string)
	data := ctx.Data["data"].(string)
	if err := database.DBConn.Model(&fetching.Status{}).Where("id = ?", ctx.Client.ID).Update("data", statusMessage).Error; err != nil {
		pipeshandler.ErrorResponse(ctx, localization.ErrorServer)
		return
	}

	// Send to other people
	conversationTokens, _, members, _, ok := conversation.PrepareConversationTokens(ctx)
	if !ok {
		return
	}

	for _, token := range conversationTokens {

		// Make sure the conversation is a private one
		if len(members[token.Conversation]) > 2 {
			continue
		}

		// Get the other member to send the status to
		var otherMember caching.StoredMember
		for _, member := range members[token.Conversation] {
			if member.TokenID != token.ID {
				otherMember = member
			}
		}

		// Send the status event
		caching.SendEventToMembers([]caching.StoredMember{otherMember}, StatusEvent(statusMessage, data, token.Conversation, token.ID, ""))
	}

	// Send the status to other devices
	caching.CSNode.SendClient(ctx.Client.ID, StatusEvent(statusMessage, data, "", ctx.Client.ID, ":o"))

	pipeshandler.SuccessResponse(ctx)
}

func StatusEvent(st string, data string, conversation string, ownToken string, suffix string) pipes.Event {
	return pipes.Event{
		Name: "acc_st" + suffix,
		Data: map[string]interface{}{
			"c":  conversation,
			"o":  ownToken,
			"st": st,
			"d":  data,
		},
	}
}
