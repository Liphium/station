package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type getDateRequest struct {
	Id string `json:"id"` // Id of the friend in the vault
}

// Route: /account/friends/get_receive_date
func getReceiveDate(c *fiber.Ctx) error {

	// Parse the request
	var req getDateRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}

	// Get the friendship from the database
	var friendship database.Friendship
	if err := database.DBConn.Where("id = ? AND account = ?", req.Id, accId).Take(&friendship).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.JSON(map[string]interface{}{
		"success": true,
		"date":    friendship.LastPacket,
	})
}

type updateDateRequest struct {
	Id   string `json:"id"`   // Id of the friend in the vault
	Date string `json:"date"` // Time of the last packet (encrypted)
}

// Route: /account/friends/update_receive_date
func updateReceiveDate(c *fiber.Ctx) error {

	// Parse the request
	var req updateDateRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}

	// Make sure someone doesn't store their whole house in here
	if len(req.Date) >= 200 {
		return integration.InvalidRequest(c, "invalid date")
	}

	// Get the friendship from the database
	if err := database.DBConn.Model(&database.Friendship{}).Where("id = ? AND account = ?", req.Id, accId).
		Update("last_packet", req.Date).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
