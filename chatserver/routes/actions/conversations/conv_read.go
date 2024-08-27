package conversation_actions

import (
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type GenericTokenConfirmAction struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// Action: conv_read
func HandleRead(c *fiber.Ctx, action GenericTokenConfirmAction) error {

	// Validate the token
	token, err := caching.ValidateToken(action.ID, action.Token)
	if err != nil {
		return integration.InvalidRequest(c, "token is invalid")
	}

	// Update read state
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("id = ?", token.ID).Update("last_read", time.Now().UnixMilli()).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
