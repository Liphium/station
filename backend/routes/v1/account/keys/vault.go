package keys

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/vault/get
func getVaultKey(c *fiber.Ctx) error {

	// Get account
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Get public key
	var key database.VaultKey
	if database.DBConn.Where("id = ?", accId).Take(&key).Error != nil {
		return util.FailedRequest(c, localization.ErrorKeyNotFound, nil)
	}

	if key.Key == "" {
		return util.FailedRequest(c, localization.ErrorKeyNotFound, nil)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"key":     key.Key,
	})
}

// Route: /account/keys/vault/set
func setVaultKey(c *fiber.Ctx) error {

	var req setRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get account
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}

	var acc database.Account
	if database.DBConn.Where("id = ?", accId).Take(&acc).Error != nil {
		return util.InvalidRequest(c)
	}

	if database.DBConn.Where("id = ?", accId).Take(&database.VaultKey{}).Error == nil {
		return util.FailedRequest(c, localization.ErrorKeyAlreadySet, nil)
	}

	// Set vault key
	if database.DBConn.Create(&database.VaultKey{
		ID:  accId,
		Key: req.Key,
	}).Error != nil {
		return util.FailedRequest(c, localization.ErrorServer, nil)
	}

	return util.SuccessfulRequest(c)
}
