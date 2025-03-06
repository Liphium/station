package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type addFriendRequest struct {
	Payload     string `json:"payload"` // Encrypted payload
	ReceiveDate string `json:"receive_date"`
}

// Route: /account/friends/add
func addFriend(c *fiber.Ctx) error {

	// Parse request
	var req addFriendRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Make sure the date isn't garbage
	if len(req.ReceiveDate) >= 150 {
		return util.InvalidRequest(c)
	}

	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Get the latest version
	var version int64
	if err := database.DBConn.Model(&database.VaultEntry{}).Select("max(version)").Where("account = ?", accId).Scan(&version).Error; err != nil {
		version = 0
	}

	// Check if the account has too many friends
	var friendCount int64
	if err := database.DBConn.Model(&database.Friendship{}).Where("account = ?", accId).Count(&friendCount).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	if friendCount >= MaximumFriends {
		return util.FailedRequest(c, localization.ErrorFriendLimitReached(MaximumFriends), nil)
	}

	// Create friendship
	friendship := database.Friendship{
		ID:         auth.GenerateToken(12),
		Account:    accId,
		Payload:    req.Payload,
		LastPacket: req.ReceiveDate,
		Version:    version + 1,
	}
	if err := database.DBConn.Create(&friendship).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"id":      friendship.ID,
		"version": version + 1,
	})
}
