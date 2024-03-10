package status

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/nodes"

	"github.com/gofiber/fiber/v2"
)

type offlineRequest struct {
	Token string `json:"token"`
}

func offline(c *fiber.Ctx) error {

	// Parse request
	var req offlineRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get node
	var requested node.Node
	if err := database.DBConn.Where("token = ?", req.Token).Take(&requested).Error; err != nil {
		return util.InvalidRequest(c)
	}

	// Update status
	nodes.TurnOff(&requested, node.StatusStopped)

	if err := database.DBConn.Save(&requested).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.SuccessfulRequest(c)
}
