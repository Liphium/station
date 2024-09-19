package vault

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
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

	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}
	if err := database.DBConn.Model(&properties.VaultEntry{}).Where("id = ? AND account = ?", req.Entry, accId).Update("payload", req.Payload).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
