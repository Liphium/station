package message_actions

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: msg_list_after
func HandleListAfter(c *fiber.Ctx, token conversations.ConversationToken, after uint64) error {

	// Get the messages
	var messages []conversations.Message
	if err := database.DBConn.Order("creation ASC").Where("conversation = ? AND creation > ?", token.Conversation, after).Limit(12).Find(&messages).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return the messages
	return integration.ReturnJSON(c, fiber.Map{
		"success":  true,
		"messages": messages,
	})
}
