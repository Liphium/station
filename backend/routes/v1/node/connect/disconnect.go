package connect

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type disconnectRequest struct {
	Node      uint   `json:"node"`
	NodeToken string `json:"token"`
	Session   string `json:"session"`
}

// Route: /node/disconnect
func Disconnect(c *fiber.Ctx) error {

	// Parse request
	var req disconnectRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Check node
	_, err := nodes.Node(req.Node, req.NodeToken)
	if err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Disconnect account
	if database.DBConn.Model(&database.Session{}).Where("id = ?", req.Session).Update("node", 0).Error != nil {
		util.Log.Println("Failed to disconnect account", req.Session)
		return integration.FailedRequest(c, localization.ErrorServer, nil)
	}

	return integration.SuccessfulRequest(c)
}
