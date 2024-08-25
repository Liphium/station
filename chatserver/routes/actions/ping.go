package remote_action_routes

import (
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type pingAction struct {
	Echo string `json:"echo"`
}

// Route: /actions/ping
func pingTest(c *fiber.Ctx) error {

	// Parse the action
	var req remoteActionRequest[pingAction]
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "action wasn't valid")
	}

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"message": "Hello from a different Liphium town!",
		"echo":    req.Data.Echo,
	})
}
