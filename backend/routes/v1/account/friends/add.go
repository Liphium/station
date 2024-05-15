package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/gofiber/fiber/v2"
)

type addFriendRequest struct {
	Hash        string `json:"hash"`    // Payload hash
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

	// Make sure the date aren't garbage
	if len(req.ReceiveDate) >= 150 {
		return util.InvalidRequest(c)
	}

	// Check if the account has too many friends
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}
	var friendCount int64
	if err := database.DBConn.Model(&properties.Friendship{}).Where("account = ?", accId).Count(&friendCount).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	if friendCount >= MaximumFriends {
		return util.FailedRequest(c, "limit.reached", nil)
	}

	// Check if it already exists
	if database.DBConn.Model(&properties.Friendship{}).Where("account = ? AND hash = ?", accId, req.Hash).Take(&properties.Friendship{}).Error == nil {
		return util.FailedRequest(c, "already.exists", nil)
	}

	// Create friendship
	friendship := properties.Friendship{
		ID:         auth.GenerateToken(12),
		Account:    accId,
		Hash:       req.Hash,
		Payload:    req.Payload,
		LastPacket: req.ReceiveDate,
	}
	if err := database.DBConn.Create(&friendship).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"id":      friendship.ID,
		"hash":    friendship.Hash,
	})
}
