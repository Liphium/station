package files

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

type infoRequest struct {
	Id string `json:"id"`
}

// Route: /account/files_unauth/info
func fileInfo(c *fiber.Ctx) error {

	var req infoRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get file info
	var cloudFile database.CloudFile
	if err := database.DBConn.Select("id,name,size,account").Where("id = ?", req.Id).Take(&cloudFile).Error; err != nil {
		return util.InvalidRequest(c)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"file":    cloudFile,
	})
}
