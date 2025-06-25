package profile

import (
	"errors"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/main/integration"
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
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get account (to notify about name & tag changes)
	var acc database.Account
	if err := database.DBConn.Where("id = ?", req.ID).Take(&acc).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get profile (to update profile picture, description, ...)
	var profile database.Profile = database.Profile{}
	if err := database.DBConn.Where("id = ?", req.ID).Take(&profile).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"profile":      profile,
		"name":         acc.Username,
		"display_name": acc.DisplayName,
	})
}
