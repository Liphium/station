package files

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type infoRequest struct {
	Id string `json:"id"`
}

// Route: /account/files_unauth/info
func fileInfo(c *fiber.Ctx) error {

	var req infoRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get file info
	var cloudFile database.CloudFile
	if err := database.DBConn.Select("id,name,size,account").Where("id = ?", req.Id).Take(&cloudFile).Error; err != nil {
		return integration.InvalidRequest(c, "database error")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"file":    cloudFile,
	})
}
