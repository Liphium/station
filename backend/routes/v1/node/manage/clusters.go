package manage

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

type clusterRequest struct {
	Token string `json:"token"`
}

func clusterList(c *fiber.Ctx) error {

	// Parse request
	var req clusterRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	var ct node.NodeCreation
	if err := database.DBConn.Where("token = ?", req.Token).Take(&ct).Error; err != nil {
		return util.InvalidRequest(c)
	}

	var clusters []node.Cluster
	if err := database.DBConn.Find(&clusters).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success":  true,
		"clusters": clusters,
	})
}
