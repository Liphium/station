package remote_action_routes

import (
	"github.com/gofiber/fiber/v2"
)

// Route: /actions/ping
func pingTest(c *fiber.Ctx, action struct {
	Echo string `json:"echo"`
}) error {

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Hello from a different Liphium town!",
		"echo":    action.Echo,
	})
}
