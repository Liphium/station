package auth_routes

import (
	login_routes "github.com/Liphium/station/backend/routes/v1/accounts/auth/login"
	register_routes "github.com/Liphium/station/backend/routes/v1/accounts/auth/register"
	sso_routes "github.com/Liphium/station/backend/routes/v1/accounts/auth/sso"
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
