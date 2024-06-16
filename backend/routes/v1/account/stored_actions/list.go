package stored_actions

import (
	"sort"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/gofiber/fiber/v2"
)

type returnableStoredAction struct {
	Id            string `json:"id"`
	Payload       string `json:"payload"`
	Authenticated bool   `json:"a"`
}

// Route: /account/stored_actions/list
func listStoredActions(c *fiber.Ctx) error {

	// Get stored actions
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}
	var returnables = []returnableStoredAction{}

	// Get all normal stored actions and add them as non authenticated ones
	var storedActions []properties.StoredAction
	if database.DBConn.Where("account = ?", accId).Find(&storedActions).Error != nil {
		return util.FailedRequest(c, "server.error", nil)
	}
	for _, storedAction := range storedActions {
		returnables = append(returnables, returnableStoredAction{
			Id:            storedAction.ID,
			Payload:       storedAction.Payload,
			Authenticated: false,
		})
	}

	// Get all authenticated stored actions and mark them as such in the returning phase
	var aStoredActions []properties.AStoredAction
	if database.DBConn.Where("account = ?", accId).Find(&aStoredActions).Error != nil {
		return util.FailedRequest(c, "server.error", nil)
	}
	for _, storedAction := range aStoredActions {
		returnables = append(returnables, returnableStoredAction{
			Id:            storedAction.ID,
			Payload:       storedAction.Payload,
			Authenticated: true,
		})
	}

	// Sort stored actions by created_at
	sort.Slice(storedActions, func(i, j int) bool {
		return storedActions[i].CreatedAt < storedActions[j].CreatedAt
	})

	// Get authenticated stored action key
	var storedActionKey account.StoredActionKey
	if database.DBConn.Where(&account.StoredActionKey{ID: accId}).Take(&storedActionKey).Error != nil {

		// Generate new stored action key
		storedActionKey = account.StoredActionKey{
			ID:  accId,
			Key: auth.GenerateToken(StoredActionTokenLength),
		}

		// Save stored action key
		if err := database.DBConn.Create(&storedActionKey).Error; err != nil {
			return util.FailedRequest(c, "server.error", err)
		}
	}

	// Return stored actions
	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"key":     storedActionKey.Key,
		"actions": returnables,
	})
}
