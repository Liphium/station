package app

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/app"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

func listApps(c *fiber.Ctx) error {

	var apps []app.App
	database.DBConn.Find(&apps)

	return util.ReturnJSON(c, apps)
}
