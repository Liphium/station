package status

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type onlineRequest struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}

func online(c *fiber.Ctx) error {

	// Parse request
	var req onlineRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get node
	requested, err := nodes.Node(req.ID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, "couldn't get node")
	}

	// Update status
	nodes.TurnOff(&requested, database.StatusStarted)

	// Send adoption
	var foundNodes []database.Node
	var startedNodes []database.NodeEntity
	if err := database.DBConn.Where(&database.Node{
		AppID:  requested.AppID,
		Status: database.StatusStarted,
	}).Find(&foundNodes).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	for _, n := range foundNodes {
		if n.ID != requested.ID {
			if err := n.SendPing(); err != nil {

				util.Log.Println("Found offline node: " + n.Domain + "! Shutting down..")

				nodes.TurnOff(&n, database.StatusStopped)
			} else {
				startedNodes = append(startedNodes, n.ToEntity())
			}
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"nodes":   startedNodes,
	})
}
