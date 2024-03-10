package stored_actions

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/auth"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
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
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	if req.Account == "" || req.Payload == "" {
		return util.InvalidRequest(c)
	}

	// Get account
	var acc account.Account
	if err := database.DBConn.Where("id = ?", req.Account).Take(&acc).Error; err != nil {
		return util.InvalidRequest(c)
	}

	// Create stored action
	storedAction := properties.StoredAction{
		ID:      auth.GenerateToken(12),
		Account: acc.ID,
		Payload: req.Payload,
	}

	// Check if stored action is authenticated
	if req.Key != "" {

		// Check if stored action limit is reached
		var storedActionCount int64
		if err := database.DBConn.Model(&properties.AStoredAction{}).Where("account = ?", acc.ID).Count(&storedActionCount).Error; err != nil {
			return util.FailedRequest(c, "server.error", err)
		}

		if storedActionCount >= AuthenticatedStoredActionLimit {
			return util.FailedRequest(c, "limit.reached", nil)
		}

		var storedActionKey account.StoredActionKey
		if err := database.DBConn.Where(&account.StoredActionKey{ID: req.Account}).Take(&storedActionKey).Error; err != nil {
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
			return util.FailedRequest(c, "server.error", err)
		}

	} else {

		// Check if stored action limit is reached
		var storedActionCount int64
		if err := database.DBConn.Model(&properties.StoredAction{}).Where("account = ?", acc.ID).Count(&storedActionCount).Error; err != nil {
			return util.FailedRequest(c, "server.error", err)
		}

		if storedActionCount >= StoredActionLimit {
			return util.FailedRequest(c, "limit.reached", nil)
		}

		// Save stored action
		if err := database.DBConn.Create(&storedAction).Error; err != nil {
			return util.FailedRequest(c, "server.error", err)
		}
	}

	// Send to node if possible
	sendStoredActionTo(acc.ID, req.Key != "", storedAction)

	// Return success
	return util.SuccessfulRequest(c)
}

func sendStoredActionTo(accId string, authenticated bool, storedAction properties.StoredAction) {

	var session account.Session
	if err := database.DBConn.Where("account = ? AND node != ?", accId, 0).Take(&session).Error; err == nil {

		// No error handling, cause it doesn't matter if it couldn't send
		requests.SendEventToNode(session.Node, accId, requests.Event{
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
