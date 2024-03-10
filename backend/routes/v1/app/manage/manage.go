package manage

import "github.com/gofiber/fiber/v2"

func Setup(router fiber.Router) {
	router.Post("/add", addApp)
	router.Post("/remove", removeApp)
}
