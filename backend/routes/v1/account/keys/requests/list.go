package key_request_routes

import (
	"errors"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Route: /account/keys/requests/list
func list(c *fiber.Ctx) error {

	// Get the account
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Get all key requests for account
	var requests []properties.KeyRequest = []properties.KeyRequest{}
	if err := database.DBConn.Where("account = ?", accId).Find(&requests).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// Return the requests as JSON
	return util.ReturnJSON(c, fiber.Map{
		"success":  true,
		"requests": requests,
	})
}
