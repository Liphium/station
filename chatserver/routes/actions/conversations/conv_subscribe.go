package conversation_actions

import (
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/handler/account"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/gofiber/fiber/v2"
)

type RemoteSubscribeAction struct {
	Tokens []database.SentConversationToken `json:"tokens"`
	Status string                           `json:"status"`
	Data   string                           `json:"data"`
	Node   string                           `json:"node"`
}

// Action: conv_sub
func HandleRemoteSubscription(c *fiber.Ctx, action RemoteSubscribeAction) error {

	// Make sure decentralization is enabled
	if !action_helpers.IsDecentralizationEnabled() {
		return integration.FailedRequest(c, localization.ErrorDecentralizationDisabled, nil)
	}

	// Check if there are too many tokens
	if len(action.Tokens) > 500 {
		return integration.InvalidRequest(c, "too many tokens")
	}

	// Validate the tokens
	conversationTokens, missingTokens, conversationIds, err := caching.ValidateTokens(&action.Tokens)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get the conversation info
	info, err := GetConversationInfo(conversationTokens)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Add adapters for remote subscription to conversations
	for _, token := range conversationTokens {
		if token.Activated {
			caching.CSNode.AdaptWS(pipes.Adapter{
				ID: "s-" + token.ID,
				Receive: func(ctx *pipes.Context) error {
					// Send the event to the token through a remote event channel
					_, err := integration.PostRequestTC(action.Node, "/event_channel/send", fiber.Map{
						"id":    token.ID,
						"token": token.Token,
						"event": *ctx.Event,
					})
					return err
				},

				// Remove the adapter if there is an error
				OnError: func(err error) {
					caching.CSNode.RemoveAdapterWS("s-" + token.Token)
				},
			})
		}
	}

	// Send the status to everyone in a goroutine
	go func() {

		// Wait a little bit for this because the other node could still be processing the response
		time.Sleep(time.Second * 2)

		// Grab all the members
		members, err := caching.LoadMembersArray(conversationIds)
		if err != nil {
			return
		}

		// Send the status to all the conversations
		for _, token := range conversationTokens {

			// Make sure it's only send to private conversations
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
			caching.SendEventToMembers([]caching.StoredMember{otherMember}, account.StatusEvent(action.Status, action.Data, token.Conversation, token.ID, ""))
		}
	}()

	// Synchronize messages for the local tokens
	go func() {
		// Wait for the client to receive the response
		time.Sleep(1 * time.Second)

		// Go through local tokens to add them to the message sync queue (if desired)
		for _, token := range conversationTokens {
			if token.Activated && token.LastSync != -1 {
				if err := caching.AddSyncToQueue(caching.SyncData{
					TokenID:      token.ID,
					Conversation: token.Conversation,
					Since:        token.LastSync,
				}); err != nil {
					util.Log.Println("error completing message sync for ", token.ID, ":", err)
				}
			}
		}
	}()

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"info":    info,
		"missing": missingTokens,
		"node":    util.OwnPath,
	})
}

// Returned to give all the information about a conversation the client needs
type ConversationInfo struct {
	Version int64  `json:"v"`
	Reads   string `json:"r"`
}

// Returns an array of conversation info
func GetConversationInfo(tokens []database.ConversationToken) (map[string]ConversationInfo, error) {
	convInfo := make(map[string]ConversationInfo, len(tokens))
	for _, token := range tokens {

		// Get the version of the conversation
		var version int64
		if err := database.DBConn.Model(&database.Conversation{}).Select("version").Where("id = ?", token.Conversation).Take(&version).Error; err != nil {
			return nil, err
		}

		// Set conversation info
		convInfo[token.Conversation] = ConversationInfo{
			Version: version,
			Reads:   token.Reads,
		}
	}

	return convInfo, nil
}
