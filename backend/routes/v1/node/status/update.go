package status

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type updateRequest struct {
	Token     string `json:"token"`
	NewStatus uint   `json:"newStatus"`
}

func update(c *fiber.Ctx) error {

	// Parse request
	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get node
	var requested database.Node
	database.DBConn.Where("token = ?", req.Token).Take(&requested)

	if requested.ID == 0 {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Update status
	requested.Status = req.NewStatus
	database.DBConn.Save(&requested)

	return integration.SuccessfulRequest(c)
}
