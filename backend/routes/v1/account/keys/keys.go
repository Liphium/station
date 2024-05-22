package keys

import (
	key_request_routes "github.com/Liphium/station/backend/routes/v1/account/keys/requests"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {

	// Routes for the public key
	router.Post("/public/get", getPublicKey)
	router.Post("/public/set", setPublicKey)

	// Routes to get and set the profile key
	router.Post("/profile/get", getProfileKey)
	router.Post("/profile/set", setProfileKey)

	// Routes to get and set the signature public key
	router.Post("/signature/get", getSignatureKey)
	router.Post("/signature/set", setSignatureKey)

	// Routes to perform a key synchronization request
	router.Route("/requests", key_request_routes.SetupRoutes)
}
