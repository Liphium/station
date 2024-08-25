package remote_action_routes

import (
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type connectionActivateAction struct {
	ID string `json:"id"`
}

// Route: /actions/conv_activate
func activateConversation(c *fiber.Ctx) error {

	// Parse the action with the request generic
	var req remoteActionRequest[connectionActivateAction]
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "action wasn't valid")
	}

	return integration.SuccessfulRequest(c)
}
