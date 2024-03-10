package manage

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/auth"

	"github.com/gofiber/fiber/v2"
)

type regenRequest struct {
	Node uint `json:"node"` // Node ID
}

func regenToken(c *fiber.Ctx) error {

	// Parse body to remove request
	var req regenRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check permission
	if !util.Permission(c, util.PermissionAdmin) {
		return util.InvalidRequest(c)
	}

	if req.Node == 0 {
		return util.FailedRequest(c, "invalid", nil)
	}

	var node node.Node
	if err := database.DBConn.Take(&node, req.Node).Error; err != nil {
		return util.FailedRequest(c, "not.found", err)
	}

	node.Token = auth.GenerateToken(300)

	if err := database.DBConn.Save(&node).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.SuccessfulRequest(c)
}
