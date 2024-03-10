package pipeshroutes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(router fiber.Router, shouldDoSocketless bool) {
	router.Route("/gateway", gatewayRouter)  // gatewayRouter is a function in gateway.go
	router.Route("/connect", adoptionRouter) // adoption is a function in adoption.go

	if shouldDoSocketless {
		router.Post("/adoption/socketless", socketless) // socketless is a function in socketless.go
	}
}
