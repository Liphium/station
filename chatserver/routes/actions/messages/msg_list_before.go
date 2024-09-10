package message_actions

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: msg_list_before
func HandleListBefore(c *fiber.Ctx, token conversations.ConversationToken, before uint64) error {

	// Get the messages
	var messages []conversations.Message
	if err := database.DBConn.Order("creation DESC").Where("conversation = ? AND creation < ?", token.Conversation, before).Limit(12).Find(&messages).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send the messages
	return integration.ReturnJSON(c, fiber.Map{
		"success":  true,
		"messages": messages,
	})
}
