package message_actions

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: msg_get
func HandleGet(c *fiber.Ctx, token database.ConversationToken, messageId string) error {

	// Get message
	var message database.Message
	if err := database.DBConn.Where("id = ? AND conversation = ?", messageId, token.Conversation).Take(&message).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return the message
	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"message": message,
	})
}
