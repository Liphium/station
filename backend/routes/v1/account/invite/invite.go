package invite_routes

import "github.com/gofiber/fiber/v2"

func Authorized(router fiber.Router) {
	router.Post("/generate", generateInvite)
	router.Post("/get_all", getAllInformation)
}
