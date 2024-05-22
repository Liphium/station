package key_request_routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(router fiber.Router) {
	router.Post("/check", check)
}
