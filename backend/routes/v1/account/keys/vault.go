package keys

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/vault/get
func getVaultKey(c *fiber.Ctx) error {

	// Get account
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Get public key
	var key account.VaultKey
	if database.DBConn.Where("id = ?", accId).Take(&key).Error != nil {
		return util.FailedRequest(c, "not.found", nil)
	}

	if key.Key == "" {
		return util.FailedRequest(c, "not.found", nil)
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
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	var acc account.Account
	if database.DBConn.Where("id = ?", accId).Take(&acc).Error != nil {
		return util.InvalidRequest(c)
	}

	if database.DBConn.Where("id = ?", accId).Take(&account.VaultKey{}).Error == nil {
		return util.FailedRequest(c, "already.set", nil)
	}

	// Set vault key
	if database.DBConn.Create(&account.VaultKey{
		ID:  accId,
		Key: req.Key,
	}).Error != nil {
		return util.FailedRequest(c, "server.error", nil)
	}

	return util.SuccessfulRequest(c)
}
