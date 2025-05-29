package files

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type listRequest struct {
	Page int `json:"page"`
}

// Route: /account/files/list
func listFiles(c *fiber.Ctx) error {

	var req listRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Check if the page is valid
	if req.Page < 0 {
		return integration.InvalidRequest(c, "invalid page")
	}

	accId := verify.InfoLocals(c).GetAccount()

	// Get files
	var files []database.CloudFile
	if database.DBConn.Where("account = ?", accId).Order("created_at").Offset(20*req.Page).Limit(20).Find(&files).Error != nil {
		return integration.FailedRequest(c, localization.ErrorServer, nil)
	}

	// Count files to calculate amount of pages and stuff (on the client)
	var count int64
	if database.DBConn.Model(&database.CloudFile{}).Where("account = ?", accId).Count(&count).Error != nil {
		return integration.FailedRequest(c, localization.ErrorServer, nil)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"files":   files,
		"count":   count,
	})
}
