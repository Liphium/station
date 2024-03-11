package node

import (
	"github.com/Liphium/station/backend/routes/v1/node/connect"
	"github.com/Liphium/station/backend/routes/v1/node/manage"
	"github.com/Liphium/station/backend/routes/v1/node/status"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/status", status.Setup)
	router.Route("/manage", manage.Unauthorized)
	router.Post("/this", this)
	router.Post("/disconnect", connect.Disconnect)
	router.Post("/get_lowest", connect.GetLowest)

	if util.Testing {
		router.Post("/test", sendToNode)
	}
}

func Authorized(router fiber.Router) {
	router.Route("/manage", manage.Authorized)
	router.Post("/connect", connect.Connect)
	router.Post("/token", generateToken)
}
