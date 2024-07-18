package routes

import (
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
	"github.com/gofiber/fiber/v2"
)

type LeaveRoomRequest struct {
	Connection string `json:"conn"`
}

// Route: /leave
func leaveRoom(c *fiber.Ctx) error {

	var req LeaveRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request body, err: "+err.Error())
	}

	connections := caching.SSInstance.GetSessions(req.Connection)
	if len(connections) == 0 {
		util.Log.Println("couldn't leave room: token not found")
		return c.JSON(fiber.Map{
			"success": true,
		})
	}

	// Disconnect the client
	for _, conn := range connections {
		caching.SSInstance.Disconnect(req.Connection, conn)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
