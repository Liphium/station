package key_request_routes

import "github.com/gofiber/fiber/v2"

func Unauthorized(router fiber.Router) {
	router.Post("/check", check)
	router.Post("/exists", exists)
}

func Authorized(router fiber.Router) {
	router.Post("/list", list)
	router.Post("/respond", respond)
}
