package connect

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
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
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check node
	_, err := nodes.Node(req.Node, req.NodeToken)
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Disconnect account
	if database.DBConn.Model(&database.Session{}).Where("id = ?", req.Session).Update("node", 0).Error != nil {
		util.Log.Println("Failed to disconnect account", req.Session)
		return util.FailedRequest(c, localization.ErrorServer, nil)
	}

	return util.SuccessfulRequest(c)
}
