package conversation

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

type conversationInfo struct {
	Version           int64 `json:"v"`
	ReadDate          int64 `json:"r"`
	NotificationCount int64 `json:"n"`
}

// Action: conv_sub
func subscribe(ctx pipeshandler.Context) {

	if ctx.ValidateForm("tokens", "status") {
		pipeshandler.ErrorResponse(ctx, localization.InvalidRequest)
		return
	}

	conversationTokens, tokenIds, members, missingTokens, ok := PrepareConversationTokens(ctx)
	if !ok {
		pipeshandler.ErrorResponse(ctx, localization.InvalidRequest)
		return
	}

	// Update all node IDs
	if database.DBConn.Model(&conversations.ConversationToken{}).Where("id IN ?", tokenIds).Update("node", util.NodeTo64(caching.CSNode.ID)).Error != nil {
		pipeshandler.ErrorResponse(ctx, localization.ErrorServer)
		return
	}

	statusctx := ctx.Data["status"].(string)
	convInfo := make(map[string]conversationInfo, len(conversationTokens))
	adapters := make([]string, len(conversationTokens))
	for _, token := range conversationTokens {

		// Register adapter for the subscription
		caching.CSNode.AdaptWS(pipes.Adapter{
			ID: "s-" + token.Token,
			Receive: func(context *pipes.Context) error {
				client := *ctx.Client
				util.Log.Println(context.Adapter.ID, token.Token, client.ID)
				err := caching.CSNode.SendClient(ctx.Client.ID, *context.Event)
				if err != nil {
					util.Log.Println("COULDN'T SEND:", err.Error())
				}
				return err
			},
		})
		util.Log.Println("SUB", "s-"+token.Token)
		adapters = append(adapters, "s-"+token.Token)

		var memberIds []string
		var memberNodes []string
		if len(members[token.Conversation]) == 2 {
			for _, member := range members[token.Conversation] {
				if member.Token != token.Token {
					memberIds = append(memberIds, "s-"+member.Token)
					memberNodes = append(memberNodes, util.Node64(member.Node))
				}
			}
		}

		// Send the subscription event
		caching.CSNode.Pipe(pipes.ProtocolWS, pipes.Message{
			Channel: pipes.Conversation(memberIds, memberNodes),
			Event: pipes.Event{
				Name: "acc_st",
				Data: map[string]interface{}{
					"st": statusctx,
					"d":  "",
				},
			},
		})

		// Get the notification count of the current conversation
		var notificationCount int64
		if err := database.DBConn.Model(&conversations.Message{}).Where("conversation = ? AND creation > ?", token.Conversation, token.LastRead).
			Count(&notificationCount).Error; err != nil {

			// Return an error
			pipeshandler.ErrorResponse(ctx, localization.ErrorServer)
			return
		}

		// Get the version of the conversation
		var version int64
		if err := database.DBConn.Model(&conversations.Conversation{}).Select("version").Where("id = ?", token.Conversation).Take(&version).Error; err != nil {

			pipeshandler.ErrorResponse(ctx, localization.ErrorServer)
			return
		}

		// Set conversation info
		convInfo[token.Conversation] = conversationInfo{
			ReadDate:          token.LastRead,
			NotificationCount: notificationCount,
		}
	}

	// Insert adapters into cache (to be deleted when disconnecting)
	caching.InsertAdapters(ctx.Client.ID, adapters)

	pipeshandler.NormalResponse(ctx, map[string]interface{}{
		"success": true,
		"info":    convInfo,
		"missing": missingTokens,
	})
}

// Returns: conversationTokens, tokenIds, members, missingTokens, success (bool)
func PrepareConversationTokens(ctx pipeshandler.Context) ([]conversations.ConversationToken, []string, map[string][]caching.StoredMember, []string, bool) {

	tokensUnparsed := ctx.Data["tokens"].([]interface{})
	tokens := make([]conversations.SentConversationToken, len(tokensUnparsed))
	for i, token := range tokensUnparsed {
		unparsed := token.(map[string]interface{})
		tokens[i] = conversations.SentConversationToken{
			ID:    unparsed["id"].(string),
			Token: unparsed["token"].(string),
		}
	}

	if len(tokens) > 500 {
		pipeshandler.ErrorResponse(ctx, localization.InvalidRequest)
		return nil, nil, nil, nil, false
	}

	var conversationTokens []conversations.ConversationToken
	var missingTokens []string
	var err error
	conversationTokens, missingTokens, err = caching.ValidateTokens(&tokens)

	if err != nil {
		pipeshandler.ErrorResponse(ctx, localization.ErrorServer)
		return nil, nil, nil, nil, false
	}

	tokenIds := make([]string, len(conversationTokens))
	conversationIds := make([]string, len(conversationTokens))
	for i, token := range conversationTokens {
		tokenIds[i] = token.ID
		conversationIds[i] = token.Conversation
	}

	members, err := caching.LoadMembersArray(conversationIds)
	if err != nil {
		pipeshandler.ErrorResponse(ctx, localization.ErrorServer)
		return nil, nil, nil, nil, false
	}

	for id, token := range members {
		util.Log.Printf("%s %d", id, len(token))
	}

	return conversationTokens, tokenIds, members, missingTokens, true
}
