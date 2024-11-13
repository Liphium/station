package keys

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/encrypted
// ! deprecated
func getAllEncryptedKeys(c *fiber.Ctx) error {
	// TODO: Remove in a future protocol upgrade

	// Get account
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
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
