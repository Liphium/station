package node

import (
	node_action_routes "github.com/Liphium/station/backend/routes/v1/node/actions"
	"github.com/Liphium/station/backend/routes/v1/node/connect"
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/actions", node_action_routes.Unauthorized)

	router.Post("/this", this)
	router.Post("/disconnect", connect.Disconnect)
	router.Post("/get_bool_setting", getBoolSetting)
	router.Post("/get_int_setting", getIntSetting)
	router.Post("/get_session", getSession)
}

func Authorized(router fiber.Router) {
	router.Post("/connect", connect.Connect)
	router.Post("/token", generateToken)
}
