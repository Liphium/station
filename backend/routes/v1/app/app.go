package app

import (
	"node-backend/routes/v1/app/manage"

	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Post("/version", getVersion)
}

func Authorized(router fiber.Router) {
	router.Route("/manage", manage.Setup)
	router.Post("/list", listApps)
}
