package townhall_accounts

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /townhall/accounts/list
func listAccounts(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Page  int    `json:"page"`
		Query string `json:"query"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get accounts
	var accounts []database.Account
	var count int64
	if req.Query != "" {
		query := "%" + req.Query + "%"

		// Search for the specified account in the database
		if database.DBConn.Order("created_at").Where("username LIKE ? OR id LIKE ? OR display_name LIKE ? OR email LIKE ?", query, query, query, query).Offset(20*req.Page).Limit(20).Find(&accounts).Error != nil {
			return util.FailedRequest(c, localization.ErrorServer, nil)
		}

		// Count accounts to calculate amount of pages and stuff (on the client)
		if database.DBConn.Model(&database.Account{}).Where("username LIKE ? OR id LIKE ? OR display_name LIKE ? OR email LIKE ?", query, query, query, query).Count(&count).Error != nil {
			return util.FailedRequest(c, localization.ErrorServer, nil)
		}
	} else {

		// Just take all accounts
		if database.DBConn.Order("created_at").Offset(20*req.Page).Limit(20).Find(&accounts).Error != nil {
			return util.FailedRequest(c, localization.ErrorServer, nil)
		}

		// Count accounts to calculate amount of pages and stuff (on the client)
		if database.DBConn.Model(&database.Account{}).Count(&count).Error != nil {
			return util.FailedRequest(c, localization.ErrorServer, nil)
		}
	}

	return util.ReturnJSON(c, fiber.Map{
		"success":  true,
		"accounts": accounts,
		"count":    count,
	})
}
