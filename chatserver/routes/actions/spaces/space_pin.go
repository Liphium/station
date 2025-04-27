package space_actions

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

// Action: space_pin
func HandleSpacePin(c *fiber.Ctx, token database.ConversationToken, action struct {
	Id         string `json:"id"`
	Underlying string `json:"underlying"`
}) error {

	// Set the underlying id of the Space
	caching.AddUnderlyingToSharedSpace(token.Conversation, action.Id, action.Underlying)

	return integration.SuccessfulRequest(c)
}
