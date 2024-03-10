package files

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

type infoRequest struct {
	Id string `json:"id"`
}

// Route: /account/files/info
func fileInfo(c *fiber.Ctx) error {

	var req infoRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get file info
	var cloudFile account.CloudFile
	if err := database.DBConn.Select("id,name,size,account,path").Where("id = ?", req.Id).Take(&cloudFile).Error; err != nil {
		return util.InvalidRequest(c)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"file":    cloudFile,
	})
}
