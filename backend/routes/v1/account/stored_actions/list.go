package stored_actions

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/auth"
	"sort"

	"github.com/gofiber/fiber/v2"
)

// Route: /account/stored_actions/list
func listStoredActions(c *fiber.Ctx) error {

	// Get stored actions
	accId := util.GetAcc(c)
	var storedActions []properties.StoredAction
	if database.DBConn.Where("account = ?", accId).Find(&storedActions).Error != nil {
		return util.FailedRequest(c, "server.error", nil)
	}

	var aStoredActions []properties.AStoredAction
	if database.DBConn.Where("account = ?", accId).Find(&aStoredActions).Error != nil {
		return util.FailedRequest(c, "server.error", nil)
	}
	for _, aStoredAction := range aStoredActions {
		storedActions = append(storedActions, properties.StoredAction(aStoredAction))
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
		"actions": storedActions,
	})
}
