package files

import (
	"strings"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/files/download/:id
func downloadFile(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		return util.InvalidRequest(c)
	}

	// Check for potentially malicious requests
	if strings.Contains(id, "/") {
		return util.InvalidRequest(c)
	}

	// Get the file from the database
	var file account.CloudFile
	if err := database.DBConn.Where("id = ?", id).Take(&file).Error; err != nil {
		return util.FailedRequest(c, "file.not_found", err)
	}

	// Send the file (it's encrypted so there is no checking of permissions required)
	return c.SendFile(saveLocation+id, true)
}
