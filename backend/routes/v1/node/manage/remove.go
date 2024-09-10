package manage

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/node"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
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
		return util.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	// Delete node
	if err := database.DBConn.Delete(node.Node{}, req.Node).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, err)
	}

	return util.SuccessfulRequest(c)
}
