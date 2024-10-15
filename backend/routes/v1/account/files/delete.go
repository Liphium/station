package files

import (
	"strings"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type deleteRequest struct {
	Id string `json:"id"`
}

// Route: /account/files/delete
func deleteFile(c *fiber.Ctx) error {

	if disabled {
		return util.FailedRequest(c, localization.ErrorFileDisabled, nil)
	}

	var req deleteRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}
	accId := verify.InfoLocals(c).GetAccount()

	// Get file
	var file database.CloudFile
	if database.DBConn.Where("account = ? AND id = ?", accId, req.Id).First(&file).Error != nil {
		return util.FailedRequest(c, localization.ErrorFileNotFound, nil)
	}

	// Check for potential malicious requests
	if strings.Contains(req.Id, "/") {
		return util.InvalidRequest(c)
	}

	// Delete the file
	if err := Delete([]string{req.Id}); err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
