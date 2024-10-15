package vault

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
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
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get friends list
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}
	var entries []database.VaultEntry
	if err := database.DBConn.Where("account = ? AND tag = ? AND updated_at > ?", accId, req.Tag, req.After).Find(&entries).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return friends list
	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"entries": entries,
	})
}
