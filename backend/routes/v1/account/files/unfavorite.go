package files

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

type unfavoriteRequest struct {
	Id string `json:"id"`
}

// Route: /account/files/unfavorite
func unfavoriteFile(c *fiber.Ctx) error {

	var req unfavoriteRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Get file
	var file account.CloudFile
	if database.DBConn.Where("account = ? AND id = ?", accId, req.Id).First(&file).Error != nil {
		return util.FailedRequest(c, "file.not_found", nil)
	}

	// Toggle favorite
	if err := database.DBConn.Model(&account.CloudFile{}).Where("account = ? AND id = ?", accId, req.Id).
		Update("favorite", false).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.SuccessfulRequest(c)
}
