package space_actions

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: space_pin_status
func HandleSpacePinStatusChange(c *fiber.Ctx, token database.ConversationToken, action struct {
	Id         string `json:"id"`
	Underlying string `json:"underlying"`
}) error {

	// Set the underlying id of the Space
	if err := caching.ChangeSpaceUnderlying(token.Conversation, action.Id, action.Underlying); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
