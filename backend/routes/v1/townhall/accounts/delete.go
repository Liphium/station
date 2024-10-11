package townhall_accounts

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/routes/v1/account/files"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /townhall/accounts/delete
func deleteAccount(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Account string `json:"account"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Delete all the data related to the account
	database.DBConn.Where("account = ?", req.Account).Delete(&database.Session{})
	database.DBConn.Where("creator = ?", req.Account).Delete(&database.Invite{})
	database.DBConn.Where("account = ?", req.Account).Delete(&database.InviteCount{})
	database.DBConn.Where("id = ?", req.Account).Delete(&database.ProfileKey{})
	database.DBConn.Where("id = ?", req.Account).Delete(&database.StoredActionKey{})
	database.DBConn.Where("id = ?", req.Account).Delete(&database.PublicKey{})
	database.DBConn.Where("id = ?", req.Account).Delete(&database.VaultKey{})
	database.DBConn.Where("id = ?", req.Account).Delete(&database.SignatureKey{})
	database.DBConn.Where("account = ?", req.Account).Delete(&database.AStoredAction{})
	database.DBConn.Where("account = ?", req.Account).Delete(&database.StoredAction{})
	database.DBConn.Where("account = ?", req.Account).Delete(&database.Friendship{})
	database.DBConn.Where("account = ?", req.Account).Delete(&database.VaultEntry{})
	database.DBConn.Where("id = ?", req.Account).Delete(&database.Profile{})
	database.DBConn.Where("account = ?", req.Account).Delete(&database.Authentication{})
	database.DBConn.Where("account = ?", req.Account).Delete(&database.KeyRequest{})
	database.DBConn.Where("id = ?", req.Account).Delete(&database.Account{})

	// Get all the files and delete all of them
	var ids []string
	if err := database.DBConn.Model(&database.CloudFile{}).Select("id").Where("account = ?", req.Account).Scan(&ids).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}
	if err := files.Delete(ids); err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
