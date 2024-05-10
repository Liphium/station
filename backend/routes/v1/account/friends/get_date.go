package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

type getDateRequest struct {
	Id string `json:"id"` // Id of the friend in the vault
}

// Route: /account/friends/get_date
func getDate(c *fiber.Ctx) error {

	// Parse the request
	var req getDateRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Get the friendship from the database
	var friendship properties.Friendship
	if err := database.DBConn.Where("id = ? AND account = ?", req.Id, accId).Take(&friendship).Error; err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	return util.ReturnJSON(c, map[string]interface{}{
		"success": true,
		"date":    friendship.LastPacket,
	})
}
