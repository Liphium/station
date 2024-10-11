package node

import (
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

func generateToken(c *fiber.Ctx) error {

	if !util.Permission(c, util.PermissionAdmin) {
		return util.InvalidRequest(c)
	}

	tk := auth.GenerateToken(200)

	// Save
	if err := database.DBConn.Create(&database.NodeCreation{
		Token: tk,
		Date:  time.Now(),
	}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"token":   tk,
	})
}
