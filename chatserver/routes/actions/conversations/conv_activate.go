package conversation_actions

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	message_routes "github.com/Liphium/station/chatserver/routes/conversations/message"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type ConnectionActivateAction struct {
	ID    string `json:"id"`    // Conversation token id
	Token string `json:"token"` // Conversation token
}

type ReturnableMember struct {
	ID   string `json:"id"`
	Rank uint   `json:"rank"`
	Data string `json:"data"`
}

// Action: conv_activate
func HandleTokenActivation(c *fiber.Ctx, action ConnectionActivateAction) error {

	// Validate the action
	if len(action.ID) == 0 || len(action.Token) == 0 {
		return integration.InvalidRequest(c, "data in action wasn't valid")
	}

	// Activate conversation
	var token conversations.ConversationToken
	if database.DBConn.Where("id = ? AND token = ?", action.ID, action.Token).First(&token).Error != nil {
		return integration.FailedRequest(c, "invalid.token", nil)
	}

	if token.Activated {
		return integration.FailedRequest(c, "already.active", nil)
	}

	// Activate token
	token.Activated = true
	token.Token = util.GenerateToken(util.ConversationTokenLength)

	if err := database.DBConn.Save(&token).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return all data
	var tokens []conversations.ConversationToken
	if err := database.DBConn.Where(&conversations.ConversationToken{Conversation: token.Conversation}).Find(&tokens).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	var members []ReturnableMember
	for _, token := range tokens {
		members = append(members, ReturnableMember{
			ID:   token.ID,
			Rank: token.Rank,
			Data: token.Data,
		})
	}

	var conversation conversations.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Take(&conversation).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	if conversation.Type == conversations.TypeGroup {

		// Increment the version by one to save the modification
		if err := action_helpers.IncrementConversationVersion(conversation); err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}

		// Send a system message to tell the group members about the new member
		err := message_routes.SendSystemMessage(token.Conversation, message_routes.GroupMemberJoin, []string{message_routes.AttachAccount(token.Data)})
		if err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
	}

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"type":    conversation.Type,
		"data":    conversation.Data,
		"token":   token.Token,
		"members": members,
	})
}
