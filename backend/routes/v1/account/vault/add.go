package vault

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/auth"

	"github.com/gofiber/fiber/v2"
)

type addEntryRequest struct {
	Tag     string `json:"tag"`     // Tag
	Payload string `json:"payload"` // Encrypted payload
}

// Route: /account/vault/add
func addEntry(c *fiber.Ctx) error {

	// Parse request
	var req addEntryRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if the account has too many entries
	accId := util.GetAcc(c)
	var entryCount int64
	if err := database.DBConn.Model(&properties.VaultEntry{}).Where("account = ?", accId).Count(&entryCount).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	if entryCount >= MaximumEntries {
		return util.FailedRequest(c, "limit.reached", nil)
	}

	// Create vault entry
	vaultEntry := properties.VaultEntry{
		ID:      auth.GenerateToken(12),
		Account: accId,
		Tag:     req.Tag,
		Payload: req.Payload,
	}
	if err := database.DBConn.Create(&vaultEntry).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"id":      vaultEntry.ID,
	})
}
