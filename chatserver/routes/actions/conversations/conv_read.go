package conversation_actions

import (
	"time"

	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: conv_read
func HandleRead(c *fiber.Ctx, token conversations.ConversationToken, _ interface{}) error {

	// Update read state
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("id = ?", token.ID).Update("last_read", time.Now().UnixMilli()).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
