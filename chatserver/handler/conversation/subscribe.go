package conversation

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler/wshandler"
)

// Action: conv_sub
func subscribe(message wshandler.Message) {

	if message.ValidateForm("tokens", "status", "date") {
		wshandler.ErrorResponse(message, localization.InvalidRequest)
		return
	}

	date := int64(message.Data["date"].(float64))
	conversationTokens, tokenIds, members, missingTokens, ok := PrepareConversationTokensWithLookup(message, true)
	if !ok {
		wshandler.ErrorResponse(message, localization.InvalidRequest)
		return
	}

	// Update all node IDs
	if database.DBConn.Model(&conversations.ConversationToken{}).Where("id IN ?", tokenIds).Update("node", util.NodeTo64(caching.Node.ID)).Error != nil {
		wshandler.ErrorResponse(message, localization.ErrorServer)
		return
	}

	statusMessage := message.Data["status"].(string)
	readDates := make(map[string]int64, len(conversationTokens))
	adapters := make([]string, len(conversationTokens))
	for _, token := range conversationTokens {

		// Register adapter for the subscription
		caching.Node.AdaptWS(pipes.Adapter{
			ID: "s-" + token.Token,
			Receive: func(ctx *pipes.Context) error {
				client := *message.Client
				util.Log.Println(ctx.Adapter.ID, token.Token, client.ID)
				err := client.SendEvent(*ctx.Event)
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
		caching.Node.Pipe(pipes.ProtocolWS, pipes.Message{
			Channel: pipes.Conversation(memberIds, memberNodes),
			Event: pipes.Event{
				Name: "acc_st",
				Data: map[string]interface{}{
					"st": statusMessage,
					"d":  "",
				},
			},
		})

		readDates[token.Conversation] = token.LastRead
		AddConversationToken(TokenTask{"s-" + token.Token, token.Conversation, date})
	}

	// Insert adapters into cache (to be deleted when disconnecting)
	caching.InsertAdapters(message.Client.ID, adapters)

	wshandler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"read":    readDates,
		"missing": missingTokens,
	})
}

// Returns: conversationTokens, tokenIds, members, missingTokens, success (bool)
func PrepareConversationTokens(message wshandler.Message) ([]conversations.ConversationToken, []string, map[string][]caching.StoredMember, []string, bool) {
	return PrepareConversationTokensWithLookup(message, false)
}

// Returns: conversationTokens, tokenIds, members, missingTokens, success (bool)
func PrepareConversationTokensWithLookup(message wshandler.Message, lookup bool) ([]conversations.ConversationToken, []string, map[string][]caching.StoredMember, []string, bool) {

	tokensUnparsed := message.Data["tokens"].([]interface{})
	tokens := make([]conversations.SentConversationToken, len(tokensUnparsed))
	for i, token := range tokensUnparsed {
		unparsed := token.(map[string]interface{})
		tokens[i] = conversations.SentConversationToken{
			ID:    unparsed["id"].(string),
			Token: unparsed["token"].(string),
		}
	}

	if len(tokens) > 500 {
		wshandler.ErrorResponse(message, localization.InvalidRequest)
		return nil, nil, nil, nil, false
	}

	var conversationTokens []conversations.ConversationToken
	var missingTokens []string
	var err error
	if lookup {
		conversationTokens, missingTokens, err = caching.ValidateTokensLookup(&tokens)
	} else {
		conversationTokens, missingTokens, err = caching.ValidateTokens(&tokens)
	}
	if err != nil {
		wshandler.ErrorResponse(message, localization.ErrorServer)
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
		wshandler.ErrorResponse(message, localization.ErrorServer)
		return nil, nil, nil, nil, false
	}

	for id, token := range members {
		util.Log.Printf("%s %d", id, len(token))
	}

	return conversationTokens, tokenIds, members, missingTokens, true
}
