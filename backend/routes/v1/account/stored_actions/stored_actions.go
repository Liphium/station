package stored_actions

import "github.com/gofiber/fiber/v2"

// Configuration
const StoredActionLimit = 10              // Max number of stored actions per account
const AuthenticatedStoredActionLimit = 20 // Max number of authenticated stored actions per account
const StoredActionTokenLength = 32        // Length of the token used to identify stored actions

// Completely public
func Unauthorized(router fiber.Router) {
	router.Post("/send_auth", sendAuthenticatedStoredAction)
}

// Authorized with account JWT
func Authorized(router fiber.Router) {
	router.Post("/list", listStoredActions)
	router.Post("/delete", deleteStoredAction)
	router.Post("/details", getDetails)
	router.Post("/send", sendStoredAction)
}
