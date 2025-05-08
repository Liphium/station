package key_request_routes

import (
	"errors"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Route: /account/keys/requests/list
func listKeyRequests(c *fiber.Ctx) error {

	// Get the account
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Get all key requests for account
	var requests []database.KeyRequest = []database.KeyRequest{}
	if err := database.DBConn.Where("account = ?", accId).Find(&requests).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Check if they are still valid
	var validRequests = []database.KeyRequest{}
	for _, request := range requests {

		// Check if the session still exists
		if err := database.DBConn.Where("account = ? AND id = ?", request.Account, request.Session).Take(&database.Session{}).Error; err != nil {

			// Delete the request if the session doesn't exist anymore
			if errors.Is(err, gorm.ErrRecordNotFound) {
				database.DBConn.Delete(&request)
				continue
			}
		}

		// Append to the final list
		validRequests = append(validRequests, request)
	}

	// Return the requests as JSON
	return util.ReturnJSON(c, fiber.Map{
		"success":  true,
		"requests": validRequests,
	})
}
