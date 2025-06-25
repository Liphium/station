package key_request_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type respondRequest struct {
	Session string `json:"session"`
	Delete  bool   `json:"delete"`
	Payload string `json:"payload"`
}

// Route: /account/keys/requests/respond
func respondToKeyRequest(c *fiber.Ctx) error {

	// Parse request
	var req respondRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get the account
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}

	// Get the account id
	sessionId, err := uuid.Parse(req.Session)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get the key synchronization request
	var request database.KeyRequest
	if err := database.DBConn.Where("session = ? AND account = ?", sessionId, accId).Take(&request).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Delete the request, if desired
	if req.Delete {
		if err := database.DBConn.Delete(&request).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}

		return integration.SuccessfulRequest(c)
	}

	// Otherwise respond to the request
	request.Payload = req.Payload
	if err := database.DBConn.Save(&request).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
