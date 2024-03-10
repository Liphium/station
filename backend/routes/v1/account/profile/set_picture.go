package profile

import (
	"log"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type setProfileRequest struct {
	Container string `json:"container"`
	File      string `json:"file"`
	Data      string `json:"data"`
}

var fileTypes = []string{
	"png",
	"jpg",
	"jpeg",
}

// Route: /account/profile/set_picture
func setProfilePicture(c *fiber.Ctx) error {

	var req setProfileRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}
	accId := util.GetAcc(c)

	var file account.CloudFile
	if err := database.DBConn.Where("id = ?", req.File).Take(&file).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	args := strings.Split(file.Id, ".")
	extension := args[len(args)-1]
	found := false
	log.Println("file extension: " + extension)
	for _, fType := range fileTypes {
		if extension == fType {
			found = true
		}
	}

	if !found {
		return util.InvalidRequest(c)
	}

	var profile properties.Profile
	err := database.DBConn.Where("id = ?", accId).Take(&profile).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return util.FailedRequest(c, "server.error", err)
	}

	if err == nil {

		// Make previous profile picture no longer saved
		if err := database.DBConn.Model(&account.CloudFile{}).Where("id = ?", profile.Picture).Update("system", false).Error; err != nil {
			return util.FailedRequest(c, "server.error", err)
		}
	}

	profile = properties.Profile{
		ID:          accId,
		Picture:     req.File,
		Container:   req.Container,
		PictureData: req.Data,
		Data:        "",
	}

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
