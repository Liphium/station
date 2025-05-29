package message_actions

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: msg_list_after
func HandleListAfter(c *fiber.Ctx, token database.ConversationToken, action struct {
	After int64  `json:"after"`
	Extra string `json:"extra"` // Extra identifier for squares
}) error {

	// Get the messages
	var messages []database.Message
	if err := database.DBConn.Order("creation ASC, id").Where("conversation = ? AND creation > ?", database.WithExtra(token.Conversation, action.Extra), action.After).Limit(30).Find(&messages).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return the messages
	return c.JSON(fiber.Map{
		"success":  true,
		"messages": messages,
	})
}
