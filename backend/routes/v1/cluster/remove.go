package cluster

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

type removeRequest struct {
	Token string `json:"token"`
	ID    uint   `json:"id"`
}

func removeCluster(c *fiber.Ctx) error {

	// Parse request
	var req removeRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	if !util.Permission(c, util.PermissionAdmin) {
		return util.InvalidRequest(c)
	}

	// Check if cluster exists
	var cluster node.Cluster
	err := database.DBConn.First(cluster, req.ID).Error

	if err == nil {
		return util.FailedRequest(c, "cluster.exists", nil)
	}

	// Remove cluster
	err = database.DBConn.Delete(cluster).Error

	if err != nil {
		return util.FailedRequest(c, "invalid", err)
	}

	return util.SuccessfulRequest(c)
}
