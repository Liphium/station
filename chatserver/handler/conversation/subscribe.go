package conversation

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database/conversations"
	conversation_actions "github.com/Liphium/station/chatserver/routes/actions/conversations"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

type subscribeAction struct {
	Tokens []conversations.SentConversationToken `json:"tokens"`
	Status string                                `json:"status"`
}

type conversationInfo struct {
	Version           int64 `json:"v"`
	ReadDate          int64 `json:"r"`
	NotificationCount int64 `json:"n"`
}

// Action: conv_sub
func subscribe(ctx *pipeshandler.Context, action subscribeAction) pipes.Event {

	// Validate the tokens
	conversationTokens, missingTokens, _, err := caching.ValidateTokens(&action.Tokens)
	if err != nil {
		return pipeshandler.ErrorResponse(ctx, localization.ErrorServer, err)
	}

	// Get the conversation info
	convInfo, err := conversation_actions.GetConversationInfo(conversationTokens)
	if err != nil {
		return pipeshandler.ErrorResponse(ctx, localization.ErrorServer, err)
	}

	// Register all the adapters
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
	}

	return pipeshandler.NormalResponse(ctx, map[string]interface{}{
		"success": true,
		"info":    convInfo,
		"missing": missingTokens,
	})
}
