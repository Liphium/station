package manage

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

type removeRequest struct {
	Node uint `json:"node"` // Node ID
}

func removeNode(c *fiber.Ctx) error {

	// Parse body to remove request
	var req removeRequest
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

	// Delete node
	if err := database.DBConn.Delete(node.Node{}, req.Node).Error; err != nil {
		return util.FailedRequest(c, "invalid", err)
	}

	return util.SuccessfulRequest(c)
}
