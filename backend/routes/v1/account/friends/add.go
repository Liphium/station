package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/main/localization"
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

	// Make sure the date isn't garbage
	if len(req.ReceiveDate) >= 150 {
		return util.InvalidRequest(c)
	}

	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Check if the friend already exists (and return id and stuff if he does)
	var friendship properties.Friendship
	if database.DBConn.Where("account = ? AND hash = ?", accId, req.Hash).Take(&friendship).Error == nil {
		return util.ReturnJSON(c, fiber.Map{
			"success": true,
			"id":      friendship.ID,
			"hash":    friendship.Hash,
		})
	}

	// Check if the account has too many friends
	var friendCount int64
	if err := database.DBConn.Model(&properties.Friendship{}).Where("account = ?", accId).Count(&friendCount).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	if friendCount >= MaximumFriends {
		return util.FailedRequest(c, localization.ErrorFriendLimitReached(MaximumFriends), nil)
	}

	// Create friendship
	friendship = properties.Friendship{
		ID:         auth.GenerateToken(12),
		Account:    accId,
		Hash:       req.Hash,
		Payload:    req.Payload,
		LastPacket: req.ReceiveDate,
	}
	if err := database.DBConn.Create(&friendship).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"id":      friendship.ID,
		"hash":    friendship.Hash,
	})
}
