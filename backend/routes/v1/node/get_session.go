package node

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/nodes"
	"github.com/Liphium/station/main/integration"
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
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Validate node token
	_, err := nodes.Node(nodeToU(req.ID), req.Token)
	if err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get the session and check if it is valid
	var session database.Session
	if err := database.DBConn.Where("id = ?", req.Session).Take(&session).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorSessionNotFound, err)
	}

	if !session.Verified {
		return integration.FailedRequest(c, localization.ErrorSessionNotVerified, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"account": session.Account.String(),
		"level":   session.PermissionLevel,
	})
}
