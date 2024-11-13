package stored_actions

import (
	"sort"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
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
	if err := database.DBConn.Where("account = ?", accId).Find(&storedActions).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
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
	if err := database.DBConn.Where("account = ?", accId).Find(&aStoredActions).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
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
	if err := database.DBConn.Where(&database.StoredActionKey{ID: accId}).Take(&storedActionKey).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return stored actions
	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"key":     storedActionKey.Key, // TODO: Deprecated, remove with next protocol version
		"actions": returnables,
	})
}
