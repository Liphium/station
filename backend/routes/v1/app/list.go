package app

import (
	"node-backend/database"
	"node-backend/entities/app"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

func listApps(c *fiber.Ctx) error {

	var apps []app.App
	database.DBConn.Find(&apps)

	return util.ReturnJSON(c, apps)
}
