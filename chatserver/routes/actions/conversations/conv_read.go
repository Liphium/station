package conversation_actions

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: conv_read
func HandleRead(c *fiber.Ctx, token database.ConversationToken, stamp int64) error {

	// Update read state
	if err := database.DBConn.Model(&database.ConversationToken{}).Where("id = ?", token.ID).Update("last_read", stamp).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
