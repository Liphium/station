package session

import (
	"node-backend/routes/v1/session/manage"

	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	router.Route("/manage", manage.Setup)
	router.Post("/logout", logOut)
}
