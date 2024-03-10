package status

import (
	"log"
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/nodes"

	"github.com/gofiber/fiber/v2"
)

type onlineRequest struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}

func online(c *fiber.Ctx) error {

	// Parse request
	var req onlineRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get node
	requested, err := nodes.Node(req.ID, req.Token)
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Update status
	nodes.TurnOff(&requested, node.StatusStarted)

	// Send adoption
	var foundNodes []node.Node
	var startedNodes []node.NodeEntity
	if err := database.DBConn.Where(&node.Node{
		AppID:  requested.AppID,
		Status: node.StatusStarted,
	}).Find(&foundNodes).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	for _, n := range foundNodes {
		if n.ID != requested.ID {
			if err := n.SendPing(); err != nil {

				log.Println("Found offline node: " + n.Domain + "! Shutting down..")

				nodes.TurnOff(&n, node.StatusStopped)
			} else {
				startedNodes = append(startedNodes, n.ToEntity())
			}
		}
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"nodes":   startedNodes,
	})
}
