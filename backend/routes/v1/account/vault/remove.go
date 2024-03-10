package vault

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"

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
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if entry exists
	accId := util.GetAcc(c)
	var entry properties.VaultEntry
	if err := database.DBConn.Where("id = ? AND account = ?", req.ID, accId).Take(&entry).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return util.FailedRequest(c, "not.found", nil)
		}

		return util.FailedRequest(c, "server.error", err)
	}

	// Delete entry
	if err := database.DBConn.Delete(&entry).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.SuccessfulRequest(c)
}
