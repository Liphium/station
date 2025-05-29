package space_actions

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

// Action: space_add
func HandleSpaceAddition(c *fiber.Ctx, token database.ConversationToken, action struct {
	Server       string `json:"server"`
	Id           string `json:"id"`
	UnderlyingId string `json:"underlying_id"`
	Name         string `json:"name"`
	Container    string `json:"container"` // Space connection container
}) error {

	// Try to add the Space
	if exists, msg := caching.StoreSharedSpace(
		token.Conversation,
		action.Server,
		action.Id,
		action.Name,
		action.UnderlyingId,
		action.Container,
	); msg != nil || exists {

		// If it already exists, send a different response
		if exists {
			return c.JSON(fiber.Map{
				"success": true,
				"exists":  true,
			})
		}

		// If it doesn't exist yet, send the message
		return integration.FailedRequest(c, msg, nil)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"exists":  false,
	})
}
