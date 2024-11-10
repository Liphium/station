package vault

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/backend/util/verify"
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
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if the account has too many entries
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}
	var entryCount int64
	if err := database.DBConn.Model(&database.VaultEntry{}).Where("account = ?", accId).Count(&entryCount).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	if entryCount >= MaximumEntries {
		return util.FailedRequest(c, localization.ErrorVaultLimitReached(MaximumEntries), nil)
	}

	// Create vault entry
	vaultEntry := database.VaultEntry{
		ID:      auth.GenerateToken(12),
		Account: accId,
		Tag:     req.Tag,
		Payload: req.Payload,
	}
	if err := database.DBConn.Create(&vaultEntry).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"id":      vaultEntry.ID,
	})
}
