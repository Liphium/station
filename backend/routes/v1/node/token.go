package node

import (
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

func generateToken(c *fiber.Ctx) error {

	if !verify.InfoLocals(c).HasPermission(verify.PermissionAdmin) {
		return integration.InvalidRequest(c, "invalid request")
	}

	tk := auth.GenerateToken(200)

	// Save
	if err := database.DBConn.Create(&database.NodeCreation{
		Token: tk,
		Date:  time.Now(),
	}).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"token":   tk,
	})
}
