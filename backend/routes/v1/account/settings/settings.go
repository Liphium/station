package settings_routes

import "github.com/gofiber/fiber/v2"

// Authorized routes
func Authorized(router fiber.Router) {
	router.Post("/change_name", changeName)
	router.Post("/change_password", changePassword)
}
