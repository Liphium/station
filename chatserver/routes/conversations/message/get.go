package message_routes

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type messageGetRequest struct {
	TokenID string `json:"token_id"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

// Route: /conversations/message/get
func get(c *fiber.Ctx) error {

	// Parse the request
	var req messageGetRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "request not valid")
	}

	// Validate the token
	token, err := caching.ValidateToken(req.TokenID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, "token not valid")
	}

	// Get message
	var message conversations.Message
	if err := database.DBConn.Where("id = ? AND conversation = ?", req.Message, token.Conversation).Take(&message).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"message": message,
	})
}
