package account

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
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
	if err := util.BodyParser(c, &req); err != nil {
		util.Log.Println(err)
		return util.InvalidRequest(c)
	}

	// Get account
	var acc database.Account
	if err := database.DBConn.Select("id", "username", "display_name").Where("username = ?", req.Name).Take(&acc).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	var pub database.PublicKey
	if err := database.DBConn.Select("key").Where("id = ?", acc.ID).Take(&pub).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	var signaturePub database.SignatureKey
	if err := database.DBConn.Select("key").Where("id = ?", acc.ID).Take(&signaturePub).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success":      true,
		"id":           acc.ID,
		"name":         acc.Username,
		"display_name": acc.DisplayName,
		"sg":           signaturePub.Key,
		"pub":          pub.Key,
	})
}
