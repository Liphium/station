package recovery_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/recovery/get
func getRecoveryData(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		SessionToken  string `json:"session_token"`
		RecoveryToken string `json:"recovery_token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get the session
	var session database.Session
	if err := database.DBConn.Where("token = ?", req.SessionToken).Take(&session).Error; err != nil {
		return integration.InvalidRequest(c, "invalid session")
	}

	// Check if the recovery token is valid
	var recoveryToken database.RecoveryToken
	if err := database.DBConn.Where("account = ? AND token = ?", session.Account, req.RecoveryToken).Take(&recoveryToken).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorInvalidRecoveryToken, err)
	}

	// Return the payload to the client
	return c.JSON(fiber.Map{
		"success": true,
		"data":    recoveryToken.Data,
	})
}
