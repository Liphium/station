package profile

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/profile/remove_picture
func removeProfilePicture(c *fiber.Ctx) error {

	// Get the account id
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Get the current profile
	var profile properties.Profile = properties.Profile{}
	err := database.DBConn.Where("id = ?", accId).Take(&profile).Error

	// Check if the profile was found (error has to be gorm.ErrRecordNotFound here cause excluded before)
	if err == nil {

		// TODO: Delete the file in the future
		// Make previous profile picture no longer saved when it wasn't found
		if err := database.DBConn.Model(&account.CloudFile{}).Where("id = ?", profile.Picture).Update("system", false).Error; err != nil {
			return util.FailedRequest(c, "server.error", err)
		}

		profile.Picture = ""
		profile.PictureData = ""

		// Save new profile
		if err := database.DBConn.Save(&profile).Error; err != nil {
			return util.FailedRequest(c, "server.error", err)
		}
	}

	return util.SuccessfulRequest(c)
}
