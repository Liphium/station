package keys

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/signature/get
func getSignatureKey(c *fiber.Ctx) error {

	// Get account
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}

	// Get public key
	var key database.SignatureKey
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

// Route: /account/keys/signature/set
func setSignatureKey(c *fiber.Ctx) error {

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

	if database.DBConn.Where("id = ?", accId).Take(&database.SignatureKey{}).Error == nil {
		return integration.FailedRequest(c, localization.ErrorKeyAlreadySet, nil)
	}

	// Set public key
	if database.DBConn.Create(&database.SignatureKey{
		ID:  accId,
		Key: req.Key,
	}).Error != nil {
		return integration.FailedRequest(c, localization.ErrorServer, nil)
	}

	return integration.SuccessfulRequest(c)
}
