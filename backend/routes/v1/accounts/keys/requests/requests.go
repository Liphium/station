package key_request_routes

import "github.com/gofiber/fiber/v2"

func Unauthorized(router fiber.Router) {
	router.Post("/check", checkKeyRequest)
	router.Post("/exists", doesKeyRequestExist)
}

func Authorized(router fiber.Router) {
	router.Post("/list", listKeyRequests)
	router.Post("/respond", respondToKeyRequest)
}
