package conversation_actions

import (
	"log"
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/handler/account"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/gofiber/fiber/v2"
)

type RemoteSubscribeAction struct {
	Tokens   []conversations.SentConversationToken `json:"tokens"`
	Status   string                                `json:"status"`
	Data     string                                `json:"data"`
	SyncDate int64                                 `json:"sync"` // Time of last sent message for message sync
	Node     string                                `json:"node"`
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
		if action.SyncDate != -1 {
			for _, token := range conversationTokens {
				if token.Activated {
					if err := caching.AddSyncToQueue(caching.SyncData{
						TokenID:      token.ID,
						Conversation: token.Conversation,
						Since:        action.SyncDate,
					}); err != nil {
						util.Log.Println("error completing message sync for ", token.ID, ":", err)
					}
				}
			}
		}
	}()
	log.Println("missing tokens: ", missingTokens)

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"info":    info,
		"missing": missingTokens,
		"node":    util.OwnPath,
	})
}

// Returned to give all the information about a conversation the client needs
type ConversationInfo struct {
	Version           int64 `json:"v"`
	ReadDate          int64 `json:"r"`
	NotificationCount int64 `json:"n"`
}

// Returns an array of conversation info
func GetConversationInfo(tokens []conversations.ConversationToken) (map[string]ConversationInfo, error) {
	convInfo := make(map[string]ConversationInfo, len(tokens))
	for _, token := range tokens {

		// Get the notification count of the current conversation
		var notificationCount int64
		if err := database.DBConn.Model(&conversations.Message{}).Where("conversation = ? AND creation > ?", token.Conversation, token.LastRead).
			Count(&notificationCount).Error; err != nil {
			return nil, err
		}

		// Get the version of the conversation
		var version int64
		if err := database.DBConn.Model(&conversations.Conversation{}).Select("version").Where("id = ?", token.Conversation).Take(&version).Error; err != nil {
			return nil, err
		}

		// Set conversation info
		convInfo[token.Conversation] = ConversationInfo{
			Version:           version,
			ReadDate:          token.LastRead,
			NotificationCount: notificationCount,
		}
	}

	return convInfo, nil
}
