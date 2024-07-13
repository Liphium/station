package message_routes

import (
	"log"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type messageListBeforeRequest struct {
	TokenID string `json:"token_id"`
	Token   string `json:"token"`
	Before  int64  `json:"before"`
}

// Route: /conversations/message/list_before
func listBefore(c *fiber.Ctx) error {

	// Parse the request
	var req messageListBeforeRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "request not valid")
	}

	// Validate the token
	token, err := caching.ValidateToken(req.TokenID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, "token not valid")
	}

	// Get the messages
	var messages []conversations.Message
	if err := database.DBConn.Order("creation DESC").Where("conversation = ? AND creation < ?", token.Conversation, req.Before).Limit(12).Find(&messages).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	log.Println(len(messages))

	return integration.ReturnJSON(c, fiber.Map{
		"success":  true,
		"messages": messages,
	})
}
