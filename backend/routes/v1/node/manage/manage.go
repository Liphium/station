package manage

import (
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Post("/new", newNode)
	router.Post("/clusters", clusterList)
}

func Authorized(router fiber.Router) {
	router.Post("/remove", removeNode)
	router.Post("/regen", regenToken)
}
