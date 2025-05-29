package stored_actions

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Route: /account/stored_actions/send_auth
func sendAuthenticatedStoredAction(c *fiber.Ctx) error {

	// Parse request
	var req sendRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	if req.Account == "" || req.Payload == "" {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Parse account id from request
	id, err := uuid.Parse(req.Account)
	if err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get account
	var acc database.Account
	if err := database.DBConn.Where("id = ?", id).Take(&acc).Error; err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Create stored action
	storedAction := database.StoredAction{
		ID:      auth.GenerateToken(12),
		Account: acc.ID,
		Payload: req.Payload,
	}

	// Check if stored action limit is reached
	var storedActionCount int64
	if err := database.DBConn.Model(&database.AStoredAction{}).Where("account = ?", acc.ID).Count(&storedActionCount).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	if storedActionCount >= AuthenticatedStoredActionLimit {
		return integration.FailedRequest(c, localization.ErrorStoredActionLimitReached(AuthenticatedStoredActionLimit), nil)
	}

	var storedActionKey database.StoredActionKey
	if err := database.DBConn.Where(&database.StoredActionKey{ID: id}).Take(&storedActionKey).Error; err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	if storedActionKey.Key != req.Key {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Save authenticated stored action
	if err := database.DBConn.Create(&database.AStoredAction{
		ID:      storedAction.ID,
		Account: storedAction.Account,
		Payload: storedAction.Payload,
	}).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send stored action to account
	sendStoredActionTo(acc.ID, true, storedAction)

	return integration.SuccessfulRequest(c)
}
