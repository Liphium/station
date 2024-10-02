package session

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/requests"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

func logOut(c *fiber.Ctx) error {

	// Get token
	sessionId, err := util.GetSession(c)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	var session database.Session
	if !requests.GetSession(sessionId, &session) {
		return util.InvalidRequest(c)
	}

	// Log out
	database.DBConn.Delete(&session)

	return util.SuccessfulRequest(c)
}
