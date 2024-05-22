package key_request_routes

import (
	"errors"

	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/chatserver/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type checkRequest struct {
	Key       string // Public key
	Signature string // Signature
}

// Route: /account/keys/requests/check
func check(c *fiber.Ctx) error {

	// Parse the request
	var req checkRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if there is an existing key sync request
	session := util.GetSession(c)
	err := database.DBConn.Where("session = ?", session).Take(&properties.KeyRequest{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return util.FailedRequest(c, "server.error", err)
	}

	// Check if the request doesn't exist yet
	if err != nil {

		// Mock response for now
		return util.ReturnJSON(c, fiber.Map{
			"success": true,
			"created": true,
		})
	}

	// Mock response for now
	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"created": false,
	})
}
