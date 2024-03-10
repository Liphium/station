package app

import (
	"node-backend/database"
	"node-backend/entities/app"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

type getVersionRequest struct {
	App uint `json:"app"`
}

// Route: /app/version
func getVersion(c *fiber.Ctx) error {

	var req getVersionRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	var app app.App
	if database.DBConn.Where("id = ?", req.App).Take(&app).Error != nil {
		return util.FailedRequest(c, "not.found", nil)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"version": app.Version,
	})
}
