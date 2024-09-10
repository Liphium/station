package settings_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util"
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
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Check username and tag
	valid, message := standards.CheckUsername(req.Username)
	if !valid {
		return util.FailedRequest(c, message, nil)
	}

	// Change username
	err := database.DBConn.Model(&account.Account{}).Where("id = ?", accId).Update("username", req.Username).Error
	if err != nil {
		return util.FailedRequest(c, localization.ErrorUsernameTaken, err)
	}

	return util.SuccessfulRequest(c)
}
