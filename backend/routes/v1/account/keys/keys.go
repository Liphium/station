package keys

import "github.com/gofiber/fiber/v2"

func Authorized(router fiber.Router) {
	router.Post("/public/get", getPublicKey)
	router.Post("/public/set", setPublicKey)

	router.Post("/profile/get", getProfileKey)
	router.Post("/profile/set", setProfileKey)

	router.Post("/signature/get", getSignatureKey)
	router.Post("/signature/set", setSignatureKey)
}
