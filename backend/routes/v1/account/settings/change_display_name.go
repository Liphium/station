package settings_routes

import (
	"log"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type changeDisplayNameRequest struct {
	Name string `json:"name"`
}

// Route: /account/settings/change_display_name
func changeDisplayName(c *fiber.Ctx) error {

	var req changeDisplayNameRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Make sure the name isn't weird data
	if message := standards.CheckDisplayName(req.Name); message != nil {
		return util.FailedRequest(c, message, nil)
	}

	// Get account from database
	var acc database.Account
	if err := database.DBConn.Where("id = ?", accId).Take(&acc).Error; err != nil {
		log.Println("account no found")
		return util.InvalidRequest(c)
	}

	// Update the display name in the account
	acc.DisplayName = req.Name

	// Save new profile
	if err := database.DBConn.Save(&acc).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
