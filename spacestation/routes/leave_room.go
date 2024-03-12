package routes

import (
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipeshandler"
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

	connections := pipeshandler.GetSessions(req.Connection)
	if len(connections) == 0 {
		return c.JSON(fiber.Map{
			"success": true,
		})
	}

	for _, conn := range connections {
		connection, valid := pipeshandler.Get(req.Connection, conn)
		if !valid {
			continue
		}

		connection.Conn.Close()
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
