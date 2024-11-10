package stored_actions

import (
	"sort"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
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
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}
	var returnables = []returnableStoredAction{}

	// Get all normal stored actions and add them as non authenticated ones
	var storedActions []database.StoredAction
	if database.DBConn.Where("account = ?", accId).Find(&storedActions).Error != nil {
		return util.FailedRequest(c, localization.ErrorServer, nil)
	}
	for _, storedAction := range storedActions {
		returnables = append(returnables, returnableStoredAction{
			Id:            storedAction.ID,
			Payload:       storedAction.Payload,
			Authenticated: false,
		})
	}

	// Get all authenticated stored actions and mark them as such in the returning phase
	var aStoredActions []database.AStoredAction
	if database.DBConn.Where("account = ?", accId).Find(&aStoredActions).Error != nil {
		return util.FailedRequest(c, localization.ErrorServer, nil)
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
	var storedActionKey database.StoredActionKey
	if database.DBConn.Where(&database.StoredActionKey{ID: accId}).Take(&storedActionKey).Error != nil {

		// Generate new stored action key
		storedActionKey = database.StoredActionKey{
			ID:  accId,
			Key: auth.GenerateToken(StoredActionTokenLength),
		}

		// Save stored action key
		if err := database.DBConn.Create(&storedActionKey).Error; err != nil {
			return util.FailedRequest(c, localization.ErrorServer, err)
		}
	}

	// Return stored actions
	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"key":     storedActionKey.Key,
		"actions": returnables,
	})
}
