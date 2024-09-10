package key_request_routes

import (
	"errors"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type existsRequest struct {
	Token string `json:"token"` // Session token
}

// Route: /account/keys/requests/exists
func exists(c *fiber.Ctx) error {

	// Parse the request
	var req existsRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get the session and check the token
	var session account.Session
	if err := database.DBConn.Where("token = ?", req.Token).Take(&session).Error; err != nil {
		return util.InvalidRequest(c)
	}

	// Check if there is an existing key sync request
	var keyRequest properties.KeyRequest
	err := database.DBConn.Where("session = ?", session.ID).Take(&keyRequest).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"exists":  err == nil,
	})
}
