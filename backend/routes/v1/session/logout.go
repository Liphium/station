package session

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/requests"
	"github.com/gofiber/fiber/v2"
)

func logOut(c *fiber.Ctx) error {

	// Get token
	sessionId := util.GetSession(c)

	var session account.Session
	if !requests.GetSession(sessionId, &session) {
		return util.InvalidRequest(c)
	}

	// Log out
	database.DBConn.Delete(&session)

	return util.SuccessfulRequest(c)
}
