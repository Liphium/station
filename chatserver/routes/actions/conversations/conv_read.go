package conversation_actions

import (
	"time"

	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: conv_read
func HandleRead(c *fiber.Ctx, token database.ConversationToken, _ interface{}) error {

	// Update read state
	newStamp := time.Now().UnixMilli()
	if err := database.DBConn.Model(&database.ConversationToken{}).Where("id = ?", token.ID).Update("last_read", newStamp).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"time":    newStamp,
	})
}
