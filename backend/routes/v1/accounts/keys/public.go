package keys

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/public/get
func getPublicKey(c *fiber.Ctx) error {

	// Get account
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		util.Log.Println("couldn't get account uuid", verify.InfoLocals(c))
		return integration.InvalidRequest(c, "invalid account id")
	}

	// Get public key
	var key database.PublicKey
	if database.DBConn.Where("id = ?", accId).Take(&key).Error != nil {
		return integration.FailedRequest(c, localization.ErrorKeyNotFound, nil)
	}

	if key.Key == "" {
		return integration.FailedRequest(c, localization.ErrorKeyNotFound, nil)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"key":     key.Key,
	})
}

type setRequest struct {
	Key string `json:"key"`
}

// Route: /account/keys/public/set
func setPublicKey(c *fiber.Ctx) error {

	var req setRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get account
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}

	var acc database.Account
	if database.DBConn.Where("id = ?", accId).Take(&acc).Error != nil {
		return integration.InvalidRequest(c, "invalid account")
	}

	if database.DBConn.Where("id = ?", accId).Take(&database.PublicKey{}).Error == nil {
		return integration.FailedRequest(c, localization.ErrorKeyAlreadySet, nil)
	}

	// Set public key
	if database.DBConn.Create(&database.PublicKey{
		ID:  accId,
		Key: req.Key,
	}).Error != nil {
		return integration.FailedRequest(c, localization.ErrorServer, nil)
	}

	return integration.SuccessfulRequest(c)
}
