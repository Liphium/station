package profile

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/routes/v1/account/files"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/profile/remove_picture
func removeProfilePicture(c *fiber.Ctx) error {

	// Get the account id
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Get the current profile
	var profile database.Profile = database.Profile{}
	err = database.DBConn.Where("id = ?", accId).Take(&profile).Error

	// Check if the profile was found (error has to be gorm.ErrRecordNotFound here cause excluded before)
	if err == nil {

		// Delete the file from the server (no longer needed)
		if err := files.Delete([]string{profile.Picture}); err != nil {
			return util.FailedRequest(c, localization.ErrorServer, err)
		}

		profile.Picture = ""
		profile.Container = ""
		profile.PictureData = ""

		// Save new profile
		if err := database.DBConn.Save(&profile).Error; err != nil {
			return util.FailedRequest(c, localization.ErrorServer, err)
		}
	}

	return util.SuccessfulRequest(c)
}
