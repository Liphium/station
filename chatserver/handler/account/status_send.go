package account

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/fetching"
	"github.com/Liphium/station/chatserver/handler/conversation"
	"github.com/Liphium/station/chatserver/util"
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

		var memberIds []string
		var memberNodes []string
		util.Log.Printf("%d", len(members[token.Conversation]))
		if len(members[token.Conversation]) == 2 {
			for _, member := range members[token.Conversation] {
				if member.Token != token.Token {
					memberIds = append(memberIds, "s-"+member.Token)
					memberNodes = append(memberNodes, util.Node64(member.Node))
				}
			}
		}
		util.Log.Printf("Sending to %d members", len(memberIds))

		// Send the subscription event
		caching.CSNode.Pipe(pipes.ProtocolWS, pipes.Message{
			Channel: pipes.Conversation(memberIds, memberNodes),
			Event:   StatusEvent(statusMessage, data, token.Conversation, token.ID, ""),
		})
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
