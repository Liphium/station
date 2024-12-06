package profile

import (
	"github.com/Liphium/station/backend/database"
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

		// TODO: Delete the file in the future
		// Make previous profile picture no longer saved when it wasn't found
		if err := database.DBConn.Model(&database.CloudFile{}).Where("id = ?", profile.Picture).Update("system", false).Error; err != nil {
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
