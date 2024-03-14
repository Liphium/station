package pipeshroutes

import (
	"github.com/Liphium/station/pipes"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router, local *pipes.LocalNode, shouldDoSocketless bool) {

	// Inject middleware to add the local node to all requests
	router.Use(func(c *fiber.Ctx) error {
		c.Locals("local", local)
		return c.Next()
	})

	router.Route("/gateway", func(router fiber.Router) {
		gatewayRouter(router, local)
	})
	router.Route("/connect", func(router fiber.Router) {
		adoptionRouter(router, local)
	})

	if shouldDoSocketless {
		router.Post("/adoption/socketless", socketless)
	}
}
