package post_routes

import (
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Post("/list_after", listAfter)
}

func Authorized(router fiber.Router) {
	router.Post("/create", createPost)
}
