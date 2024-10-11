package node

import (
	node_action_routes "github.com/Liphium/station/backend/routes/v1/node/actions"
	"github.com/Liphium/station/backend/routes/v1/node/connect"
	"github.com/Liphium/station/backend/routes/v1/node/manage"
	"github.com/Liphium/station/backend/routes/v1/node/status"
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/status", status.Setup)
	router.Route("/manage", manage.Unauthorized)
	router.Route("/actions", node_action_routes.Unauthorized)

	router.Post("/this", this)
	router.Post("/disconnect", connect.Disconnect)
	router.Post("/get_lowest", connect.GetLowest)
	router.Post("/get_bool_setting", getBoolSetting)
}

func Authorized(router fiber.Router) {
	router.Post("/connect", connect.Connect)
	router.Post("/token", generateToken)
}
