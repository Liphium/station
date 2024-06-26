package account

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

type getRequest struct {
	ID string `json:"id"`
}

// Route: /account/get
func getAccount(c *fiber.Ctx) error {

	// Parse request
	var req getRequest
	if err := util.BodyParser(c, &req); err != nil {
		util.Log.Println(err)
		return util.InvalidRequest(c)
	}

	// Get account
	var acc account.Account
	if err := database.DBConn.Select("username", "tag").Where("id = ?", req.ID).Take(&acc).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	var pub account.PublicKey
	if err := database.DBConn.Select("key").Where("id = ?", req.ID).Take(&pub).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	var signaturePub account.SignatureKey
	if err := database.DBConn.Select("key").Where("id = ?", req.ID).Take(&signaturePub).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"name":    acc.Username,
		"sg":      signaturePub.Key,
		"pub":     pub.Key,
	})
}
