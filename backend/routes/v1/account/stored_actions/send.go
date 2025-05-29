package stored_actions

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/backend/util/requests"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type sendRequest struct {
	Account string `json:"account"`
	Payload string `json:"payload"`
	Key     string `json:"key"` // Authentication key
}

// Route: /account/stored_actions/send
func sendStoredAction(c *fiber.Ctx) error {

	// Parse request
	var req sendRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	if req.Account == "" || req.Payload == "" {
		return integration.InvalidRequest(c, "invalid payload or acc")
	}

	// Parse account id from request
	id, err := uuid.Parse(req.Account)
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}

	// Get account
	var acc database.Account
	if err := database.DBConn.Where("id = ?", req.Account).Take(&acc).Error; err != nil {
		return integration.InvalidRequest(c, "invalid account")
	}

	// Create stored action
	storedAction := database.StoredAction{
		ID:      auth.GenerateToken(12),
		Account: acc.ID,
		Payload: req.Payload,
	}

	// Check if stored action is authenticated
	if req.Key != "" {

		// Check if stored action limit is reached
		var storedActionCount int64
		if err := database.DBConn.Model(&database.AStoredAction{}).Where("account = ?", id).Count(&storedActionCount).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}

		if storedActionCount >= AuthenticatedStoredActionLimit {
			return integration.FailedRequest(c, localization.ErrorVaultLimitReached(AuthenticatedStoredActionLimit), nil)
		}

		var storedActionKey database.StoredActionKey
		if err := database.DBConn.Where(&database.StoredActionKey{ID: id}).Take(&storedActionKey).Error; err != nil {
			return integration.InvalidRequest(c, "invalid stored action key (database)")
		}

		if storedActionKey.Key != req.Key {
			return integration.InvalidRequest(c, "invalid stored action key (check)")
		}

		// Save authenticated stored action
		if err := database.DBConn.Create(&database.AStoredAction{
			ID:      storedAction.ID,
			Account: storedAction.Account,
			Payload: storedAction.Payload,
		}).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}

	} else {

		// Check if stored action limit is reached
		var storedActionCount int64
		if err := database.DBConn.Model(&database.StoredAction{}).Where("account = ?", acc.ID).Count(&storedActionCount).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}

		if storedActionCount >= StoredActionLimit {
			return integration.FailedRequest(c, localization.ErrorStoredActionLimitReached(StoredActionLimit), nil)
		}

		// Save stored action
		if err := database.DBConn.Create(&storedAction).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
	}

	// Send to node if possible
	sendStoredActionTo(acc.ID, req.Key != "", storedAction)

	// Return success
	return integration.SuccessfulRequest(c)
}

func sendStoredActionTo(accId uuid.UUID, authenticated bool, storedAction database.StoredAction) {

	var session database.Session
	if err := database.DBConn.Where("account = ? AND node != ?", accId, 0).Take(&session).Error; err == nil {

		// No error handling, cause it doesn't matter if it couldn't send
		requests.SendEventToNode(session.Node, accId.String(), requests.Event{
			Sender: "0",
			Name:   "s_a", // Stored action
			Data: map[string]interface{}{
				"a":       authenticated, // Authenticated
				"id":      storedAction.ID,
				"payload": storedAction.Payload,
			},
		})
	}
}
