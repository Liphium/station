package vault

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

type ListRequest struct {
	After int64  `json:"after"`
	Tag   string `json:"tag"`
}

// Route: /account/vault/list
func listEntries(c *fiber.Ctx) error {

	// Parse request
	var req ListRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Get friends list
	accId := util.GetAcc(c)
	var entries []properties.VaultEntry
	if err := database.DBConn.Where("account = ? AND tag = ? AND updated_at > ?", accId, req.Tag, req.After).Find(&entries).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Return friends list
	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"entries": entries,
	})
}
