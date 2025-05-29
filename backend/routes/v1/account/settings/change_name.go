package settings_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type changeNameRequest struct {
	Username string `json:"name"`
}

// Route: /account/settings/change_name
func changeName(c *fiber.Ctx) error {

	var req changeNameRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid id")
	}

	// Check username and tag
	if msg := standards.CheckUsername(req.Username); msg != nil {
		return integration.FailedRequest(c, msg, nil)
	}

	// Change username
	err = database.DBConn.Model(&database.Account{}).Where("id = ?", accId).Update("username", req.Username).Error
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorUsernameTaken, err)
	}

	return integration.SuccessfulRequest(c)
}
