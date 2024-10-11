package status

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
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
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get node
	requested, err := nodes.Node(req.ID, req.Token)
	if err != nil {
		return util.InvalidRequest(c)
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
		return util.FailedRequest(c, localization.ErrorServer, err)
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

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"nodes":   startedNodes,
	})
}
