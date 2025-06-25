package vault

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
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
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Check if the account has too many entries
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}
	var entryCount int64
	if err := database.DBConn.Model(&database.VaultEntry{}).Where("account = ?", accId).Count(&entryCount).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	if entryCount >= MaximumEntries {
		return integration.FailedRequest(c, localization.ErrorVaultLimitReached(MaximumEntries), nil)
	}

	// Get the best version
	var version int64
	if err := database.DBConn.Model(&database.VaultEntry{}).Select("max(version)").Where("account = ? AND tag = ?", accId, req.Tag).Scan(&version); err != nil {
		version = 0
	}

	// Create vault entry
	vaultEntry := database.VaultEntry{
		ID:      auth.GenerateToken(12),
		Account: accId,
		Tag:     req.Tag,
		Payload: req.Payload,
		Version: version + 1,
	}
	if err := database.DBConn.Create(&vaultEntry).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"id":      vaultEntry.ID,
		"version": version + 1,
	})
}
