package routes

import (
	"slices"

	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
	"github.com/gofiber/fiber/v2"
)

// Route: /enc/join
func joinSpace(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Id string `json:"id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request body")
	}

	// Check if there is a room with this id
	room, valid := caching.GetRoom(req.Id)
	if !valid {
		return integration.FailedRequest(c, localization.ErrorSpaceNotFound, nil)
	}

	// Get all the connections to the room
	clientId, token, err := newConnectionToken(c, room.ID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"token":   token,
		"client":  clientId,
	})
}

// Create a new (only join) connection token to a room
func newConnectionToken(c *fiber.Ctx, roomId string) (string, string, error) {

	// Get all the connections to the room
	connections, valid := caching.GetAllAdapters(roomId)
	if !valid {
		return "", "", integration.FailedRequest(c, localization.ErrorSpaceNotFound, nil)
	}

	// Generate a random client id
	clientId := util.GenerateToken(12)
	for slices.Contains(connections, clientId) {
		clientId = util.GenerateToken(12)
	}

	// Create a connection token
	extra := "oj-" + roomId // To highlight that the user can only join
	token, err := caching.SSInstance.GenerateToken(clientId, clientId, extra, uint(util.NodeTo64(caching.SSNode.ID)))
	if err != nil {
		return "", "", integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return clientId, token, nil
}
