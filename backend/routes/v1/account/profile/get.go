package profile

import (
	"errors"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type getProfileRequest struct {
	ID string `json:"id"`
}

// Route: /account/profile/get
func getProfile(c *fiber.Ctx) error {

	var req getProfileRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get account (to notify about name & tag changes)
	var acc database.Account
	if err := database.DBConn.Where("id = ?", req.ID).Take(&acc).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get profile (to update profile picture, description, ...)
	var profile database.Profile = database.Profile{}
	if err := database.DBConn.Where("id = ?", req.ID).Take(&profile).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success":      true,
		"profile":      profile,
		"name":         acc.Username,
		"display_name": acc.DisplayName,
		"picture":      profile.Picture,
		"container":    profile.Container,
	})
}
