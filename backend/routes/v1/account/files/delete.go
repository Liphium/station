package files

import (
	"strings"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type deleteRequest struct {
	Id string `json:"id"`
}

// Route: /account/files/delete
func deleteFile(c *fiber.Ctx) error {

	if disabled {
		return integration.FailedRequest(c, localization.ErrorFileDisabled, nil)
	}

	var req deleteRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}
	accId := verify.InfoLocals(c).GetAccount()

	// Get file
	var file database.CloudFile
	if database.DBConn.Where("account = ? AND id = ?", accId, req.Id).First(&file).Error != nil {
		return integration.FailedRequest(c, localization.ErrorFileNotFound, nil)
	}

	// Check for potential malicious requests
	if strings.Contains(req.Id, "/") {
		return integration.InvalidRequest(c, "contains directory")
	}

	// Delete the file
	if err := Delete([]string{req.Id}); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
