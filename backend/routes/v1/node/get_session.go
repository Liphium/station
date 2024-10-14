package node

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /node/get_session
func getSession(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		ID      string `json:"id"`
		Token   string `json:"token"`
		Session string `json:"session"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Validate node token
	_, err := nodes.Node(nodeToU(req.ID), req.Token)
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Get the session and check if it is valid
	var session database.Session
	if err := database.DBConn.Where("id = ?", req.Session).Take(&session).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorSessionNotFound, err)
	}

	if !session.Verified {
		return util.FailedRequest(c, localization.ErrorSessionNotVerified, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"account": session.Account.String(),
		"level":   session.PermissionLevel,
	})
}
