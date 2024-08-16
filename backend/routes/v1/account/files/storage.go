package files

import (
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/files/storage
func getStorageUsage(c *fiber.Ctx) error {

	// Get the account id
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Count the account's storage
	storage, err := CountTotalStorage(accId)
	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Return all the data
	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"amount":  storage,
		"max":     maxTotalStorage,
	})
}
