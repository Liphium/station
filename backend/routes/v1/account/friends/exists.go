package friends

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

type existsRequest struct {
	Hash string `json:"hash"`
}

func existsFriend(c *fiber.Ctx) error {

	var req existsRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if the friendship exists
	accId := util.GetAcc(c)
	if err := database.DBConn.Where("account = ? AND hash = ?", accId, req.Hash).Take(&properties.Friendship{}).Error; err != nil {
		return util.FailedRequest(c, "not.found", nil)
	}

	return util.SuccessfulRequest(c)
}
