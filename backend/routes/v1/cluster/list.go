package cluster

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

func listClusters(c *fiber.Ctx) error {

	if !util.Permission(c, util.PermissionUseServices) {
		return util.InvalidRequest(c)
	}

	var clusters []node.Cluster
	database.DBConn.Find(&clusters)

	return util.ReturnJSON(c, fiber.Map{
		"success":  true,
		"clusters": clusters,
	})
}
