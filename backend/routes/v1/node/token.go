package node

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/auth"
	"time"

	"github.com/gofiber/fiber/v2"
)

func generateToken(c *fiber.Ctx) error {

	if !util.Permission(c, util.PermissionAdmin) {
		return util.InvalidRequest(c)
	}

	tk := auth.GenerateToken(200)

	// Save
	if err := database.DBConn.Create(&node.NodeCreation{
		Token: tk,
		Date:  time.Now(),
	}).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"token":   tk,
	})
}
