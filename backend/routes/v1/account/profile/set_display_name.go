package profile

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type setDisplayNameRequest struct {
	Name string `json:"name"`
}

// Route: /account/profile/set_display_name
func setDisplayName(c *fiber.Ctx) error {

	var req setDisplayNameRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Make sure the name isn't weird data (not tested, please work)
	if len(req.Name) > 1000 {
		return util.InvalidRequest(c)
	}

	// Get the current profile
	var profile properties.Profile = properties.Profile{}
	err := database.DBConn.Where("id = ?", accId).Take(&profile).Error

	// Only return if there was an error with the database
	if err != nil && err != gorm.ErrRecordNotFound {
		return util.FailedRequest(c, "server.error", err)
	}

	// Update the display name in the profile
	profile.ID = accId
	profile.DisplayName = req.Name

	// Save new profile
	if err := database.DBConn.Save(&profile).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.SuccessfulRequest(c)
}
