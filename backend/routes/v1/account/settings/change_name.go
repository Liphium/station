package settings_routes

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/standards"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

type changeNameRequest struct {
	Username string `json:"name"`
	Tag      string `json:"tag"`
}

// Route: /account/settings/change_name
func changeName(c *fiber.Ctx) error {

	var req changeNameRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}
	accId := util.GetAcc(c)

	// Check username and tag
	valid, message := standards.CheckUsernameAndTag(req.Username, req.Tag)
	if !valid {
		return util.FailedRequest(c, message, nil)
	}

	// Change username
	err := database.DBConn.Model(&account.Account{}).Where("id = ?", accId).Update("username", req.Username).Update("tag", req.Tag).Error
	if err != nil {
		return util.FailedRequest(c, "username.taken", err)
	}

	return util.SuccessfulRequest(c)
}
