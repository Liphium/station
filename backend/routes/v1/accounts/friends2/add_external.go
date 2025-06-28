package friends2_routes

import (
	"errors"
	"fmt"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/service"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var NodeProtocol = "http://"

// Route: /accounts/friends/add_external
func addFriendFromExternal(c *fiber.Ctx) error {
	var req struct {
		From       string `json:"from"` // Address of the sender
		To         string `json:"to"`   // Address of the target (on current town)
		ProfileKey string `json:"prf"`  // Profile key to decrypt from's profile (for to, sealed)
	}
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "request not valid")
	}

	return integration.SuccessfulRequest(c)
}

// Create a friend request from a target account id for an account. “accountLPH“ should be an address.
//
// Accepts the friend request in case there was one.
func createFriendRequest(accountLPH standards.LPHAddress, target uuid.UUID, profileKey string) error {

	// Make sure profile key isn't rediculous, encryption with libcgc can get big, so this *should* be fine
	if len(profileKey) > 500 {
		return fmt.Errorf("ridiculous profile key")
	}

	// Just for validating that the account actually exists TODO: Use LoadMultiple when finally there
	accountInfo, err := service.LoadAccount(accountLPH)
	if err != nil {
		return fmt.Errorf("couldn't load account from og: %s", err)
	}
	targetInfo, err := service.LoadAccount(standards.LiphiumAddress(target.String()))
	if err != nil {
		return fmt.Errorf("couldn't load target: %s", err)
	}

	var friendship database.Friendship
	found := true
	if err := database.DBConn.Where("target = ?", accountLPH).Take(&friendship).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("couldn't get friendship: %s", err)
		}
		found = false
	}

	// Accept when already exists, otherwise create
	if found {
		targetAddr := standards.LiphiumAddress(target.String())
		if err := database.DBConn.Create(&database.Friendship{
			Request:    false,
			Account:    accountLPH.String(),
			Target:     targetAddr.String(),
			ProfileKey: profileKey,
		}).Error; err != nil {
			return fmt.Errorf("couldn't create friendship: %s", err)
		}
		friendship.Request = false
		if err := database.DBConn.Save(&friendship).Error; err != nil {
			return fmt.Errorf("couldn't create friendship 2: %s", err)
		}

		eventForAccount := FriendEventFromAccountInfo(false, targetInfo, friendship.ProfileKey)
		eventForTarget := FriendEventFromAccountInfo(false, accountInfo, profileKey)

		// TODO: Register the adapter for account (in case decentralized)

		if err := service.Instance.SendOne(accountLPH.String(), eventForAccount); err != nil {
			return fmt.Errorf("couldn't send to account: %s", err)
		}
		if err := service.Instance.SendOne(targetAddr.String(), eventForTarget); err != nil {
			return fmt.Errorf("couldn't send to target: %s", err)
		}
	} else {
		targetAddr := standards.LiphiumAddress(target.String())
		if err := database.DBConn.Create(&database.Friendship{
			Request:    true,
			Account:    accountLPH.String(),
			Target:     targetAddr.String(),
			ProfileKey: profileKey,
		}).Error; err != nil {
			return fmt.Errorf("couldn't create friend request: %s", err)
		}

		eventForTarget := FriendEventFromAccountInfo(true, accountInfo, profileKey)
		if err := service.Instance.SendOne(targetAddr.String(), eventForTarget); err != nil {
			return fmt.Errorf("couldn't send request to target: %s", err)
		}
	}

	return nil
}
