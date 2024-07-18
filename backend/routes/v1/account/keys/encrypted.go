package keys

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/encrypted
func getAllEncryptedKeys(c *fiber.Ctx) error {

	// Get account
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Get vault key
	var vaultKey account.VaultKey
	if database.DBConn.Where("id = ?", accId).Take(&vaultKey).Error != nil {
		return util.FailedRequest(c, "not.found", nil)
	}

	// Get profile key
	var profileKey account.ProfileKey
	if database.DBConn.Where("id = ?", accId).Take(&profileKey).Error != nil {
		return util.FailedRequest(c, "not.found", nil)
	}

	if vaultKey.Key == "" || profileKey.Key == "" {
		return util.FailedRequest(c, "not.found", nil)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"vault":   vaultKey.Key,
		"profile": profileKey.Key,
	})
}
