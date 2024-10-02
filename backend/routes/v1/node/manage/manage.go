package manage

import (
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Post("/new", newNode)
}
