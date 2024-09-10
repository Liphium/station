package key_request_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type respondRequest struct {
	Session string `json:"session"`
	Delete  bool   `json:"delete"`
	Payload string `json:"payload"`
}

// Route: /account/keys/requests/respond
func respond(c *fiber.Ctx) error {

	// Parse request
	var req respondRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get the account
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Get the key synchronization request
	var request properties.KeyRequest
	if err := database.DBConn.Where("session = ? AND account = ?", req.Session, accId).Take(&request).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Delete the request, if desired
	if req.Delete {
		if err := database.DBConn.Delete(&request).Error; err != nil {
			return util.FailedRequest(c, localization.ErrorServer, err)
		}

		return util.SuccessfulRequest(c)
	}

	// Otherwise respond to the request
	request.Payload = req.Payload
	if err := database.DBConn.Save(&request).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
