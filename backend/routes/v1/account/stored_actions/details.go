package stored_actions

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

type detailsRequest struct {
	Username string `json:"username"`
}

// Route: /account/stored_actions/details
func getDetails(c *fiber.Ctx) error {

	// Parse request
	var req detailsRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get account
	var acc account.Account
	if err := database.DBConn.Where("username = ?", req.Username).Take(&acc).Error; err != nil {
		return util.FailedRequest(c, "not.found", err)
	}

	var key account.PublicKey
	if err := database.DBConn.Where("id = ?", acc.ID).Take(&key).Error; err != nil {
		return util.FailedRequest(c, "not.found", err)
	}

	var signatureKey account.SignatureKey
	if err := database.DBConn.Where("id = ?", acc.ID).Take(&signatureKey).Error; err != nil {
		return util.FailedRequest(c, "not.found", err)
	}

	// Return account details
	return util.ReturnJSON(c, fiber.Map{
		"success":      true,
		"account":      acc.ID,
		"display_name": acc.DisplayName,
		"key":          key.Key,
		"sg":           signatureKey.Key,
	})
}
