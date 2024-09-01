package remote_action_routes

import (
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type pingAction struct {
	Echo string `json:"echo"`
}

// Route: /actions/ping
func pingTest(c *fiber.Ctx, action struct {
	Echo string `json:"echo"`
}) error {

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"message": "Hello from a different Liphium town!",
		"echo":    action.Echo,
	})
}
