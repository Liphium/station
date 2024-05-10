package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

type updateDateRequest struct {
	Id   string `json:"id"`   // Id of the friend in the vault
	Date string `json:"date"` // Time of the last packet (encrypted)
}

// Route: /account/friends/update_date
func updateSendDate(c *fiber.Ctx) error {

	// Parse the request
	var req updateDateRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Make sure someone doesn't store their whole house in here
	if len(req.Date) >= 200 {
		return util.InvalidRequest(c)
	}

	// Get the friendship from the database
	if err := database.DBConn.Model(&properties.Friendship{}).Where("id = ? AND account = ?", req.Id, accId).
		Update("last_packet", req.Date).Error; err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
