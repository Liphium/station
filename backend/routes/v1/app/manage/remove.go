package manage

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/app"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

type removeRequest struct {
	ID uint `json:"id"`
}

func removeApp(c *fiber.Ctx) error {

	// Parse request
	var req removeRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	if !util.Permission(c, util.PermissionAdmin) {
		return util.InvalidRequest(c)
	}

	// Delete app
	if err := database.DBConn.Delete(&app.App{}, req.ID).Error; err != nil {
		return util.InvalidRequest(c)
	}

	// TOOD: Purge everything related to the app

	return util.SuccessfulRequest(c)
}
