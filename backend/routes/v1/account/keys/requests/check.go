package key_request_routes

import (
	"errors"
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type checkRequest struct {
	Token     string `json:"token"`     // Session token
	Key       string `json:"key"`       // Public key
	Signature string `json:"signature"` // Signature
}

// Route: /account/keys/requests/check
func checkKeyRequest(c *fiber.Ctx) error {

	// Parse the request
	var req checkRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get the session and check the token
	var session database.Session
	if err := database.DBConn.Where("token = ?", req.Token).Take(&session).Error; err != nil {
		return integration.InvalidRequest(c, "invalid session")
	}

	// Check if there is an existing key sync request
	var keyRequest database.KeyRequest = database.KeyRequest{}
	err := database.DBConn.Where("session = ?", session.ID).Take(&keyRequest).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Check if the request doesn't exist yet
	if err != nil {

		// Create a new key synchronization request
		keyRequest := database.KeyRequest{
			Session:   session.ID,
			Account:   session.Account,
			Key:       req.Key,
			Signature: req.Signature,
			Payload:   "",
			CreatedAt: time.Now().UnixMilli(),
		}

		if err := database.DBConn.Create(&keyRequest).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}

		// Tell the client that the request was created
		return c.JSON(fiber.Map{
			"success": true,
			"created": true,
		})
	}

	// Return the result of the key request
	return c.JSON(fiber.Map{
		"success": true,
		"created": false,
		"payload": keyRequest.Payload,
	})
}
