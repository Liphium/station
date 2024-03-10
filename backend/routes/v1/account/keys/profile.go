package keys

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/profile/get
func getProfileKey(c *fiber.Ctx) error {

	// Get account
	accId := util.GetAcc(c)

	// Get public key
	var key account.ProfileKey
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

// Route: /account/keys/profile/set
func setProfileKey(c *fiber.Ctx) error {

	var req setRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get account
	accId := util.GetAcc(c)

	var acc account.Account
	if database.DBConn.Where("id = ?", accId).Take(&acc).Error != nil {
		return util.InvalidRequest(c)
	}

	if database.DBConn.Where("id = ?", accId).Take(&account.ProfileKey{}).Error == nil {
		return util.FailedRequest(c, "already.set", nil)
	}

	// Set public key
	if database.DBConn.Create(&account.ProfileKey{
		ID:  accId,
		Key: req.Key,
	}).Error != nil {
		return util.FailedRequest(c, "server.error", nil)
	}

	return util.SuccessfulRequest(c)
}
