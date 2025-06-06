package files

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type changeTagRequest struct {
	Id  string `json:"id"`
	Tag string `json:"tag"`
}

// Route: /account/files/change_tag
func changeFileTag(c *fiber.Ctx) error {

	var req changeTagRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}
	accId := verify.InfoLocals(c).GetAccount()

	// Check if tag is valid
	if len(req.Tag) > 100 {
		return util.InvalidRequest(c)
	}

	// Get file
	var file database.CloudFile
	if database.DBConn.Where("account = ? AND id = ?", accId, req.Id).First(&file).Error != nil {
		return util.FailedRequest(c, localization.ErrorFileNotFound, nil)
	}

	// Change the tag
	if err := database.DBConn.Model(&database.CloudFile{}).Where("account = ? AND id = ?", accId, file.Id).Update("tag", req.Tag).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
