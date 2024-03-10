package routes

import (
	"github.com/Liphium/station/integration"
	"github.com/gofiber/fiber/v2"
)

func ping(c *fiber.Ctx) error {

	//! NEVER CHANGE THIS (important for other nodes and the backend (this is a standard!))
	return c.JSON(fiber.Map{
		"gateway": integration.NODE_ID,
		"app":     "space-node",
	})
}
