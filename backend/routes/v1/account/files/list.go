package files

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

type listRequest struct {
	Favorite bool  `json:"favorite"`
	Start    int64 `json:"last"` // Start data
}

// Route: /account/files/list
func listFiles(c *fiber.Ctx) error {

	var req listRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Get files
	var files []account.CloudFile
	if database.DBConn.Where("account = ? AND favorite = ? AND created_at < ?", accId, req.Start).Limit(40).Find(&[]account.CloudFile{}).Error != nil {
		return util.FailedRequest(c, "server.error", nil)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"file":    files,
	})
}
