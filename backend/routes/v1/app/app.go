package app

import (
	"github.com/Liphium/station/backend/routes/v1/app/manage"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {
	router.Route("/manage", manage.Setup)
	router.Post("/list", listApps)
}
