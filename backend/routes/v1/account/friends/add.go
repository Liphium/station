package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/friends/add
func addFriend(c *fiber.Ctx) error {

	// Parse request
	var req struct {
		Payload     string `json:"payload"` // Encrypted payload
		ReceiveDate string `json:"receive_date"`
	}
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Make sure the date isn't garbage
	if len(req.ReceiveDate) >= 150 {
		return integration.InvalidRequest(c, "invalid receive date")
	}

	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}

	// Get the latest version
	var version int64
	if err := database.DBConn.Model(&database.VaultEntry{}).Select("max(version)").Where("account = ?", accId).Scan(&version).Error; err != nil {
		version = 0
	}

	// Check if the account has too many friends
	var friendCount int64
	if err := database.DBConn.Model(&database.Friendship{}).Where("account = ?", accId).Count(&friendCount).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	if friendCount >= MaximumFriends {
		return integration.FailedRequest(c, localization.ErrorFriendLimitReached(MaximumFriends), nil)
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
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"id":      friendship.ID,
		"version": version + 1,
	})
}
