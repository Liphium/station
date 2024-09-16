package auth_routes

import (
	login_routes "github.com/Liphium/station/backend/routes/v1/account/auth/login"
	register_routes "github.com/Liphium/station/backend/routes/v1/account/auth/register"
	sso_routes "github.com/Liphium/station/backend/routes/v1/account/auth/sso"
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Post("/refresh", refreshSession)
	router.Post("/start", startAuth)
	router.Post("/form", getStartForm)

	// Setup all the auth routes
	router.Route("/login", login_routes.Unauthorized)
	router.Route("/register", register_routes.Unauthorized)
	router.Route("/sso", sso_routes.Unauthorized)
}

func Unencrypted(router fiber.Router) {
	router.Route("/sso", sso_routes.Unencrypted)
}
