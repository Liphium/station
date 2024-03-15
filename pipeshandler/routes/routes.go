package pipeshroutes

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router, local *pipes.LocalNode, instance *pipeshandler.Instance, shouldDoSocketless bool) {

	// Inject middleware to add the local node to all requests
	router.Use(func(c *fiber.Ctx) error {
		c.Locals("local", local)
		return c.Next()
	})

	router.Route("/gateway", func(router fiber.Router) {
		gatewayRouter(router, local, instance)
	})
	router.Route("/connect", func(router fiber.Router) {
		adoptionRouter(router, local, instance)
	})

	if shouldDoSocketless {
		router.Post("/adoption/socketless", socketless)
	}
}
