package conversation_actions

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type RemoteSubscribeAction struct {
	Tokens []conversations.SentConversationToken `json:"tokens"`
	Status string                                `json:"status"`
}

// Action: conv_sub
func HandleRemoteSubscription(c *fiber.Ctx, action RemoteSubscribeAction) error {

	// Check if there are too many tokens
	if len(action.Tokens) > 500 {
		return integration.InvalidRequest(c, "too many tokens")
	}

	// Validate the tokens
	conversationTokens, missingTokens, tokenIds, err := caching.ValidateTokens(&action.Tokens)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get the conversation info
	info, err := GetConversationInfo(conversationTokens)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	if database.DBConn.Model(&conversations.ConversationToken{}).Where("id IN ?", tokenIds).Updates(map[string]interface{}{
		"remote": true,
		"node":   util.NodeTo64(caching.CSNode.ID),
	}).Error != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"info":    info,
		"missing": missingTokens,
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
