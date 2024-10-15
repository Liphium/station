package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
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
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}
	if err := database.DBConn.Where("account = ? AND hash = ?", accId, req.Hash).Take(&database.Friendship{}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorFriendNotFound, nil)
	}

	return util.SuccessfulRequest(c)
}
