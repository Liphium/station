package vault

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type updateRequest struct {
	Entry   string `json:"entry"`
	Payload string `json:"payload"`
}

// Route: /account/vault/update
func updateVaultEntry(c *fiber.Ctx) error {

	// Parse the request
	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get the current account id
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}

	// Get the entry
	var entry database.VaultEntry
	if err := database.DBConn.Model(&database.VaultEntry{}).Where("id = ? AND account = ?", req.Entry, accId).Take(&entry).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get the latest version in the vault
	var version int64
	if err := database.DBConn.Model(&database.VaultEntry{}).Select("max(version)").Where("account = ? AND tag = ?", accId, entry.Tag).Scan(&version).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Update the entry to the newest version
	entry.Payload = req.Payload
	entry.Version = version + 1
	if err := database.DBConn.Save(&entry).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"version": version + 1,
	})
}
