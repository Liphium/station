package auth

import (
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Post("/refresh", refreshSession)
}

func Authorized(router fiber.Router) {
}
