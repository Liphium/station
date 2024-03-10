package auth

import (
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Post("/login/start", startLogin)
	router.Post("/refresh", refreshSession)
	router.Post("/register/start", registerStart)   // Step 1
	router.Post("/register/code", registerCode)     // Step 2
	router.Post("/register/finish", registerFinish) // Step 3
}

func Authorized(router fiber.Router) {
	router.Post("/login/step", loginStep)
}
