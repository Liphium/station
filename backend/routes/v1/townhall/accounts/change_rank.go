package townhall_accounts

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /townhall/accounts/change_rank
func changeRank(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Account string `json:"account"`
		Rank    uint   `json:"rank"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get the rank for later processing
	var rank database.Rank
	if err := database.DBConn.Where("id = ?", req.Rank).Take(&rank).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get the account and their current rank
	var account database.Account
	if err := database.DBConn.Where("id = ?", req.Account).Preload("Rank").Take(&account).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Check if the current user has permission to set that rank
	currentLevel := verify.InfoLocals(c).GetPermissionLevel()
	if currentLevel < int16(account.Rank.Level) {
		return util.FailedRequest(c, localization.ErrorNoPermission, nil)
	}

	// Check if the user is trying to upgrade their own rank to a higher one
	acc := verify.InfoLocals(c).GetAccount()
	if acc == req.Account && currentLevel < int16(rank.Level) {
		return util.FailedRequest(c, localization.ErrorNoPermission, nil)
	}

	// Check if the user is trying to demote the only admin left
	if account.Rank.Level >= uint(verify.Permissions[verify.PermissionAdmin]) && rank.Level < uint(verify.Permissions[verify.PermissionAdmin]) {

		// Get all the accounts with admin permissions
		var ranksWithAdminPerms []uint
		if err := database.DBConn.Model(&database.Rank{}).Select("id").Where("level >= ?", verify.Permissions[verify.PermissionAdmin]).Find(&ranksWithAdminPerms).Error; err != nil {
			return util.FailedRequest(c, localization.ErrorServer, err)
		}
		var count int64
		if err := database.DBConn.Model(&database.Account{}).Where("rank_id IN ?", ranksWithAdminPerms).Count(&count).Error; err != nil {
			return util.FailedRequest(c, localization.ErrorServer, err)
		}

		if count <= 1 {
			return util.FailedRequest(c, localization.ErrorOneAdminNeeded, nil)
		}
	}

	// Change the rank of the account
	if err := database.DBConn.Model(&database.Account{}).Where("id = ?", req.Account).Update("rank_id", req.Rank).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Change all the level in all the sessions of the account
	if err := database.DBConn.Model(&database.Session{}).Where("account = ?", req.Account).Update("permission_level", rank.Level).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
