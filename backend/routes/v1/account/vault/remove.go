package vault

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type removeRequest struct {
	ID string `json:"id"`
}

// Route: /account/vault/remove
func removeEntry(c *fiber.Ctx) error {

	// Parse request
	var req removeRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Check if entry exists
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}
	var entry database.VaultEntry
	if err := database.DBConn.Where("id = ? AND account = ?", req.ID, accId).Take(&entry).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return integration.FailedRequest(c, localization.ErrorEntryNotFound, nil)
		}

		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get the best version
	var version int64
	if err := database.DBConn.Model(&database.VaultEntry{}).Select("max(version)").Where("account = ? AND tag = ?", accId, entry.Tag).Scan(&version).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Mark the vault entry as deleted and remove all of its values
	entry.Payload = "-"
	entry.Deleted = true
	entry.Version = version + 1
	if err := database.DBConn.Save(&entry).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"tag":     entry.Tag,
		"version": version + 1,
	})
}
