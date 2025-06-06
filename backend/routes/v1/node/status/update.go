package status

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

type updateRequest struct {
	Token     string `json:"token"`
	NewStatus uint   `json:"newStatus"`
}

func update(c *fiber.Ctx) error {

	// Parse request
	var req updateRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get node
	var requested database.Node
	database.DBConn.Where("token = ?", req.Token).Take(&requested)

	if requested.ID == 0 {
		return util.InvalidRequest(c)
	}

	// Update status
	requested.Status = req.NewStatus
	database.DBConn.Save(&requested)

	return util.SuccessfulRequest(c)
}
