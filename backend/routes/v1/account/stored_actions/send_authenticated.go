package stored_actions

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Route: /account/stored_actions/send_auth
func sendAuthenticatedStoredAction(c *fiber.Ctx) error {

	// Parse request
	var req sendRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	if req.Account == "" || req.Payload == "" {
		return util.InvalidRequest(c)
	}

	// Parse account id from request
	id, err := uuid.Parse(req.Account)
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Get account
	var acc account.Account
	if err := database.DBConn.Where("id = ?", id).Take(&acc).Error; err != nil {
		return util.InvalidRequest(c)
	}

	// Create stored action
	storedAction := properties.StoredAction{
		ID:      auth.GenerateToken(12),
		Account: acc.ID,
		Payload: req.Payload,
	}

	// Check if stored action limit is reached
	var storedActionCount int64
	if err := database.DBConn.Model(&properties.AStoredAction{}).Where("account = ?", acc.ID).Count(&storedActionCount).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	if storedActionCount >= AuthenticatedStoredActionLimit {
		return util.FailedRequest(c, localization.ErrorStoredActionLimitReached(AuthenticatedStoredActionLimit), nil)
	}

	var storedActionKey account.StoredActionKey
	if err := database.DBConn.Where(&account.StoredActionKey{ID: id}).Take(&storedActionKey).Error; err != nil {
		return util.InvalidRequest(c)
	}

	if storedActionKey.Key != req.Key {
		return util.InvalidRequest(c)
	}

	// Save authenticated stored action
	if err := database.DBConn.Create(&properties.AStoredAction{
		ID:      storedAction.ID,
		Account: storedAction.Account,
		Payload: storedAction.Payload,
	}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send stored action to account
	sendStoredActionTo(acc.ID, true, storedAction)

	return util.SuccessfulRequest(c)
}
