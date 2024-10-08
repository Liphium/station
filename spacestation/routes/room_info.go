package routes

import (
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/gofiber/fiber/v2"
)

type roomInfoRequest struct {
	Room string `json:"room"`
}

// Route: /info
func roomInfo(c *fiber.Ctx) error {

	// Parse request
	var req roomInfoRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request body, err: "+err.Error())
	}

	room, validRoom := caching.GetRoom(req.Room)
	members, valid := caching.GetAllConnections(req.Room)
	if !valid || !validRoom {
		return integration.FailedRequest(c, localization.ErrorRoomNotFound, nil)
	}

	returnableMembers := make([]string, len(members))
	i := 0
	for _, member := range members {
		returnableMembers[i] = member.Data
		i++
	}

	return c.JSON(fiber.Map{
		"success": true,
		"start":   room.Start,
		"members": returnableMembers,
	})
}
