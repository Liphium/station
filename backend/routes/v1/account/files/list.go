package files

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

type listRequest struct {
	Page int `json:"page"`
}

// Route: /account/files/list
func listFiles(c *fiber.Ctx) error {

	var req listRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if the page is valid
	if req.Page < 0 {
		return util.InvalidRequest(c)
	}

	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Get files
	var files []account.CloudFile
	if database.DBConn.Where("account = ?", accId).Order("created_at").Offset(20*req.Page).Limit(20).Find(&files).Error != nil {
		return util.FailedRequest(c, "server.error", nil)
	}

	// Count files to calculate amount of pages and stuff (on the client)
	var count int64
	if database.DBConn.Model(&account.CloudFile{}).Where("account = ?", accId).Count(&count).Error != nil {
		return util.FailedRequest(c, "server.error", nil)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"files":   files,
		"count":   count,
	})
}
