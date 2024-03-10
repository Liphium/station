package stored_actions

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

type deleteRequest struct {
	ID string `json:"id"`
}

// Route: /account/stored_actions/delete
func deleteStoredAction(c *fiber.Ctx) error {

	// Parse request
	var req deleteRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Delete stored action
	accId := util.GetAcc(c)
	if err := database.DBConn.Where("account = ? AND id = ?", accId, req.ID).Delete(&properties.StoredAction{}).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}
	if err := database.DBConn.Where("account = ? AND id = ?", accId, req.ID).Delete(&properties.AStoredAction{}).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Return success
	return util.SuccessfulRequest(c)
}
