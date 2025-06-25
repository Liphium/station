package vault

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
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
	if err := c.BodyParser(&req); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get friends list
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}
	var entries []database.VaultEntry
	if err := database.DBConn.Where("account = ? AND tag = ? AND deleted = ? AND updated_at > ?", accId, req.Tag, false, req.After).Find(&entries).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return friends list
	return c.JSON(fiber.Map{
		"success": true,
		"entries": entries,
	})
}
