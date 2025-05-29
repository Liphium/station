package session

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/requests"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

func logOut(c *fiber.Ctx) error {

	// Get token
	sessionId, err := verify.InfoLocals(c).GetSessionUUID()
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	var session database.Session
	if !requests.GetSession(sessionId, &session) {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Log out
	database.DBConn.Delete(&session)

	return integration.SuccessfulRequest(c)
}
