package key_request_routes

import (
	"errors"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type existsRequest struct {
	Token string `json:"token"` // Session token
}

// Route: /account/keys/requests/exists
func doesKeyRequestExist(c *fiber.Ctx) error {

	// Parse the request
	var req existsRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get the session and check the token
	var session database.Session
	if err := database.DBConn.Where("token = ?", req.Token).Take(&session).Error; err != nil {
		return integration.InvalidRequest(c, "invalid session")
	}

	// Check if there is an existing key sync request
	var keyRequest database.KeyRequest
	err := database.DBConn.Where("session = ?", session.ID).Take(&keyRequest).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"exists":  err == nil,
	})
}
