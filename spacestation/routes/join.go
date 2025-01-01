package routes

import (
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
	"github.com/gofiber/fiber/v2"
)

// Route: /join
func joinRoom(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Id string `json:"id"`
	}
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "invalid request body")
	}

	// Check if there is a room with this id
	room, valid := caching.GetRoom(req.Id)
	if !valid {
		return integration.FailedRequest(c, localization.ErrorSpaceNotFound, nil)
	}

	// Create a connection token
	clientId := util.GenerateToken(12)
	caching.SSInstance.GenerateToken(clientId, clientId, room.ID, uint(util.NodeTo64(caching.SSNode.ID)))

	return integration.SuccessfulRequest(c)
}
