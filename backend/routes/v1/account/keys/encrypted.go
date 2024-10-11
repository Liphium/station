package keys

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
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
	var vaultKey database.VaultKey
	if database.DBConn.Where("id = ?", accId).Take(&vaultKey).Error != nil {
		return util.FailedRequest(c, localization.ErrorKeyNotFound, nil)
	}

	// Get profile key
	var profileKey database.ProfileKey
	if database.DBConn.Where("id = ?", accId).Take(&profileKey).Error != nil {
		return util.FailedRequest(c, localization.ErrorKeyNotFound, nil)
	}

	if vaultKey.Key == "" || profileKey.Key == "" {
		return util.FailedRequest(c, localization.ErrorKeyNotFound, nil)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"vault":   vaultKey.Key,
		"profile": profileKey.Key,
	})
}
