package friends2_routes

import (
	"errors"
	"fmt"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/service"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FriendRemoveExternalRequest struct {
	From string `json:"from"` // Address of the sender
	To   string `json:"to"`   // Address of the target (on current town)
}

// Route: /accounts/friends/remove_external
func removeFriendFromExternal(c *fiber.Ctx) error {
	var req FriendRemoveExternalRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "request not valid")
	}

	target := standards.LPHAddress(req.To)
	_, targetTown, ok := target.Split()
	if !ok {
		return integration.FailedRequest(c, localization.ErrorInvalidRequestContent, nil)
	}
	sender := standards.LPHAddress(req.From)
	_, senderTown, ok := sender.Split()
	if !ok {
		return integration.FailedRequest(c, localization.ErrorInvalidRequestContent, nil)
	}

	if targetTown != standards.CurrentTown() || senderTown == standards.CurrentTown() {
		return integration.FailedRequest(c, localization.ErrorInvalidRequestContent, nil)
	}

	err := deleteFriendAndRequest(standards.LPHAddress(req.From), target)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send removal event to own user
	eventForTarget := FriendRemoveEvent(sender)

	// TODO: Register the adapter for account (in case decentralized) idk stand so in add external

	if err := service.Instance.SendOne(target.String(), eventForTarget); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, fmt.Errorf("couldn't send to target: %s", err))
	}

	return integration.SuccessfulRequest(c)
}

// Delete a friend or friend request from a target account id or from the account itself. “accountLPH“ should be an address. "targetLPH" should be an address.
func deleteFriendAndRequest(accountLPH standards.LPHAddress, targetLPH standards.LPHAddress) (err error) {

	// Just for validating that the account actually exists TODO: Use LoadMultiple when finally there
	_, _, err = service.LoadAccount(accountLPH)
	if err != nil {
		return fmt.Errorf("couldn't load account from og: %s", err)
	}
	_, _, err = service.LoadAccount(targetLPH)
	if err != nil {
		return fmt.Errorf("couldn't load target: %s", err)
	}

	// Requested/Friend = sender
	var friendship1 database.Friendship
	found1 := true
	if err := database.DBConn.Where("target = ? AND account = ?", accountLPH, targetLPH).Take(&friendship1).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("couldn't get friendship1: %s", err)
		}
		found1 = false
	}

	// Requested/Friend = target
	var friendship2 database.Friendship
	found2 := true
	if err := database.DBConn.Where("target = ? AND account = ?", targetLPH, accountLPH).Take(&friendship2).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("couldn't get friendship2: %s", err)
		}
		found2 = false
	}

	if !found1 && !found2 {
		return fmt.Errorf("no friendship exists between the accounts")
	}

	// Delete friend or friend request
	if found1 {
		database.DBConn.Delete(&friendship1)
	}
	if found2 {
		database.DBConn.Delete(&friendship2)
	}

	return nil
}
