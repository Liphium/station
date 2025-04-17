package message_actions

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: msg_list_before
func HandleListBefore(c *fiber.Ctx, token database.ConversationToken, action struct {
	Extra  string `json:"extra"`
	Before int64  `json:"before"`
}) error {

	// Get the messages
	var messages []database.Message
	if err := database.DBConn.Order("creation DESC").Where("conversation = ? AND creation < ?", database.WithExtra(token.Conversation, action.Extra), action.Before).Limit(12).Find(&messages).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send the messages
	return integration.ReturnJSON(c, fiber.Map{
		"success":  true,
		"messages": messages,
	})
}
