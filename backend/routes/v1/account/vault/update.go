package vault

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

type updateRequest struct {
	Entry   string `json:"entry"`
	Payload string `json:"payload"`
}

// Route: /account/vault/update
func updateVaultEntry(c *fiber.Ctx) error {

	var req updateRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	accId := util.GetAcc(c)
	if err := database.DBConn.Model(&properties.VaultEntry{}).Where("id = ? AND account = ?", req.Entry, accId).Update("payload", req.Payload).Error; err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
