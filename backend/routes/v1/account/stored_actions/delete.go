package stored_actions

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
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
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}
	if err := database.DBConn.Where("account = ? AND id = ?", accId, req.ID).Delete(&database.StoredAction{}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}
	if err := database.DBConn.Where("account = ? AND id = ?", accId, req.ID).Delete(&database.AStoredAction{}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return success
	return util.SuccessfulRequest(c)
}
