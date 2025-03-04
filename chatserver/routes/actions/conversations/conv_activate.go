package conversation_actions

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	message_actions "github.com/Liphium/station/chatserver/routes/actions/messages"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
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
func HandleTokenActivation(c *fiber.Ctx, token conversations.ConversationToken, data interface{}) error {

	if token.Activated {
		return integration.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	// Activate token
	token.Activated = true
	token.Token = util.GenerateToken(util.ConversationTokenLength)

	if err := database.DBConn.Save(&token).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send a system message in case of a group
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
		err := message_actions.SendSystemMessage(token.Conversation, message_actions.GroupMemberJoin, []string{message_actions.AttachAccount(token.Data)})
		if err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
	}

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"token":   token.Token,
	})
}
