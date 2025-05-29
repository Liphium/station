package account

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type getByUsernameRequest struct {
	Name string `json:"name"`
}

// Route: /account/get_name
func getAccountByUsername(c *fiber.Ctx) error {

	// Parse request
	var req getByUsernameRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get account
	var acc database.Account
	if err := database.DBConn.Select("id", "username", "display_name").Where("username = ?", req.Name).Take(&acc).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	var pub database.PublicKey
	if err := database.DBConn.Select("key").Where("id = ?", acc.ID).Take(&pub).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	var signaturePub database.SignatureKey
	if err := database.DBConn.Select("key").Where("id = ?", acc.ID).Take(&signaturePub).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"id":           acc.ID,
		"name":         acc.Username,
		"display_name": acc.DisplayName,
		"sg":           signaturePub.Key,
		"pub":          pub.Key,
	})
}
