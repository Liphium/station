package settings_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type changeNameRequest struct {
	Username string `json:"name"`
}

// Route: /account/settings/change_name
func changeName(c *fiber.Ctx) error {

	var req changeNameRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Check username and tag
	if msg := standards.CheckUsername(req.Username); msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Change username
	err = database.DBConn.Model(&database.Account{}).Where("id = ?", accId).Update("username", req.Username).Error
	if err != nil {
		return util.FailedRequest(c, localization.ErrorUsernameTaken, err)
	}

	return util.SuccessfulRequest(c)
}
