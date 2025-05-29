package account

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/routes/v1/account/stored_actions"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/me
func me(c *fiber.Ctx) error {

	// Get session
	sessionId, err := verify.InfoLocals(c).GetSessionUUID()
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	var session database.Session
	if database.DBConn.Where(&database.Session{ID: sessionId}).Take(&session).Error != nil {
		return integration.InvalidRequest(c, "invalid session")
	}

	// Get account
	var acc database.Account
	if err := database.DBConn.Where(&database.Account{ID: session.Account}).Take(&acc).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get all valid permissions the account has
	perms := []string{}
	info := verify.InfoLocals(c)
	for name := range verify.Permissions {
		if info.HasPermission(name) {
			perms = append(perms, name)
		}
	}

	// Get all the ranks
	var ranks []database.Rank
	if err := database.DBConn.Find(&ranks).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get vault key
	var vaultKey database.VaultKey
	if database.DBConn.Where("id = ?", acc.ID).Take(&vaultKey).Error != nil {
		return integration.FailedRequest(c, localization.ErrorKeyNotFound, nil)
	}

	// Get profile key
	var profileKey database.ProfileKey
	if database.DBConn.Where("id = ?", acc.ID).Take(&profileKey).Error != nil {
		return integration.FailedRequest(c, localization.ErrorKeyNotFound, nil)
	}

	// Make sure all the keys are there
	if vaultKey.Key == "" || profileKey.Key == "" {
		return integration.FailedRequest(c, localization.ErrorKeyNotFound, nil)
	}

	// Get authenticated stored action key
	var storedActionKey database.StoredActionKey
	if database.DBConn.Where(&database.StoredActionKey{ID: acc.ID}).Take(&storedActionKey).Error != nil {

		// Generate new stored action key
		storedActionKey = database.StoredActionKey{
			ID:  acc.ID,
			Key: auth.GenerateToken(stored_actions.StoredActionTokenLength),
		}

		// Save stored action key
		if err := database.DBConn.Create(&storedActionKey).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
	}

	// Retrun details
	return c.JSON(fiber.Map{
		"success":     true,
		"account":     acc,
		"permissions": perms,
		"ranks":       ranks,
		"vault":       vaultKey.Key,
		"profile":     profileKey.Key,
		"actions":     storedActionKey.Key,
	})
}
