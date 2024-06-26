package settings_routes

import (
	"log"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

type changeDisplayNameRequest struct {
	Name string `json:"name"`
}

// Route: /account/settings/change_display_name
func changeDisplayName(c *fiber.Ctx) error {

	var req changeDisplayNameRequest
	if err := util.BodyParser(c, &req); err != nil {
		log.Println("requesto invalid")
		return util.InvalidRequest(c)
	}
	accId, valid := util.GetAcc(c)
	if !valid {
		log.Println("account id no found")
		return util.InvalidRequest(c)
	}

	// Make sure the name isn't weird data
	if valid, message := standards.CheckDisplayName(req.Name); !valid {
		return util.FailedRequest(c, message, nil)
	}

	// Get account from database
	var acc account.Account
	if err := database.DBConn.Where("id = ?", accId).Take(&acc).Error; err != nil {
		log.Println("account no found")
		return util.InvalidRequest(c)
	}

	// Update the display name in the account
	acc.DisplayName = req.Name

	// Save new profile
	if err := database.DBConn.Save(&acc).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.SuccessfulRequest(c)
}
