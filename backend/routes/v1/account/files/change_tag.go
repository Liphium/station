package files

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
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
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}
	accId := verify.InfoLocals(c).GetAccount()

	// Check if tag is valid
	if len(req.Tag) > 100 {
		return integration.InvalidRequest(c, "tag is invalid")
	}

	// Get file
	var file database.CloudFile
	if database.DBConn.Where("account = ? AND id = ?", accId, req.Id).First(&file).Error != nil {
		return integration.FailedRequest(c, localization.ErrorFileNotFound, nil)
	}

	// Change the tag
	if err := database.DBConn.Model(&database.CloudFile{}).Where("account = ? AND id = ?", accId, file.Id).Update("tag", req.Tag).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
