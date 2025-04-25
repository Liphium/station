package space_actions

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

// Action: space_rename
func HandleSpaceRename(c *fiber.Ctx, token database.ConversationToken, action struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}) error {

	// Rename the Space
	caching.RenameSharedSpace(
		token.Conversation,
		action.Id,
		action.Name,
	)

	return integration.SuccessfulRequest(c)
}
