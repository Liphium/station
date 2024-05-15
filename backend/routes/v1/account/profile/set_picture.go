package profile

import (
	"strings"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type setProfileRequest struct {
	Container string `json:"container"`
	File      string `json:"file"`
	Data      string `json:"data"`
}

// All accepted file types (GIF would be cool :sunglasses: (maybe future or sth))
var fileTypes = []string{
	"png",
	"jpg",
	"jpeg",
}

// Route: /account/profile/set_picture
func setProfilePicture(c *fiber.Ctx) error {

	// Parse the request
	var req setProfileRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Make sure the data isn't weird (let's hope I don't regret this, not tested btw xd)
	if len(req.File) > 1000 || len(req.Container) > 2000 || len(req.Data) > 1000 {
		return util.InvalidRequest(c)
	}

	// Get the profile picture file
	var file account.CloudFile
	if err := database.DBConn.Where("id = ?", req.File).Take(&file).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Check if the file extension is correct
	args := strings.Split(file.Id, ".")
	extension := args[len(args)-1]
	found := false
	util.Log.Println("file extension: " + extension)
	for _, fType := range fileTypes {
		if extension == fType {
			found = true
		}
	}

	if !found {
		return util.InvalidRequest(c)
	}

	// Get the current profile
	var profile properties.Profile = properties.Profile{}
	err := database.DBConn.Where("id = ?", accId).Take(&profile).Error

	// Only return if there was an error with the database (exclude not found)
	if err != nil && err != gorm.ErrRecordNotFound {
		return util.FailedRequest(c, "server.error", err)
	}

	// Check if the profile was found (error has to be gorm.ErrRecordNotFound here cause excluded before)
	if err == nil {

		// Make previous profile picture no longer saved when it wasn't found
		if err := database.DBConn.Model(&account.CloudFile{}).Where("id = ?", profile.Picture).Update("system", false).Error; err != nil {
			return util.FailedRequest(c, "server.error", err)
		}
	}

	// Update all things in the profile
	profile.ID = accId
	profile.Picture = req.File
	profile.Container = req.Container
	profile.PictureData = req.Data

	// Save new profile
	if err := database.DBConn.Save(&profile).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Mark new profile picture as system file
	if err := database.DBConn.Model(&account.CloudFile{}).Where("id = ?", req.File).Update("system", true).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.SuccessfulRequest(c)
}
