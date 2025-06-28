package friends2_routes

import (
	"fmt"
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/service"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util/requests"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/neogate"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var NodeProtocol = "http://"

// Route: /accounts/friends/add_external
func addFriendFromExternal(c *fiber.Ctx) error {
	var req struct {
		From string `json:"from"` // Address of the sender
		To   string `json:"to"`   // Address of the target (on current town)
	}
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "request not valid")
	}

	return integration.SuccessfulRequest(c)
}

// Create a friend request from a target account id for an account. “accountLPH“ should be an address.
//
// Accepts the friend request in case there was one.
func createFriendRequest(accountLPH standards.LPHAddress, target uuid.UUID) error {
	accountId, origin, valid := accountLPH.Split()
	if !valid {
		return fmt.Errorf("invalid address: %s", accountLPH)
	}
	accountUuid, err := uuid.FromBytes([]byte(accountId))
	if err != nil {
		return fmt.Errorf("invalid account id (not uuid prob): %s", err)
	}

	// Verify account in case not decentralized
	var accountInfo struct {
		Username     string
		DisplayName  string
		PublicKey    string
		SignatureKey string
	}
	if origin != standards.CurrentTown() {
		res, err := requests.PostRequest(origin, "/accounts/get", requests.Map{
			"id": accountUuid.String(),
		})
		if err != nil {
			return fmt.Errorf("couldn't verify account on origin: %s", err)
		}
		if !requests.ValueOr(res, "success", false) {
			return fmt.Errorf("account verification failed with: %s", requests.ValueOr(res, "message", "unknown error"))
		}
		if requests.ValueOr(res, "id", "-") != accountId {
			return fmt.Errorf("sender account id is invalid")
		}
		accountInfo.Username = requests.ValueOr(res, "name", "")
		accountInfo.DisplayName = requests.ValueOr(res, "display_name", "")
		accountInfo.PublicKey = requests.ValueOr(res, "pub", "")
		accountInfo.SignatureKey = requests.ValueOr(res, "sg", "")
	} else {
		var account database.Account
		if err := database.DBConn.Where("id = ?", accountUuid).Take(&account).Error; err != nil {
			return fmt.Errorf("couldn't get account from database: %s", err)
		}
		var publicKey database.PublicKey
		if err := database.DBConn.Where("id = ?", accountUuid).Take(&publicKey).Error; err != nil {
			return fmt.Errorf("couldn't get public key from db: %s", err)
		}
		var signatureKey database.SignatureKey
		if err := database.DBConn.Where("id = ?", accountUuid).Take(&signatureKey).Error; err != nil {
			return fmt.Errorf("couldn't get signature key from db: %s", err)
		}
		accountInfo.Username = account.Username
		accountInfo.DisplayName = account.DisplayName
		accountInfo.PublicKey = publicKey.Key
		accountInfo.SignatureKey = signatureKey.Key
	}

	if accountInfo.Username == "" || accountInfo.DisplayName == "" || accountInfo.PublicKey == "" || accountInfo.SignatureKey == "" {
		return fmt.Errorf("invalid account info")
	}

	// TODO: Accept when already exists
	if err := database.DBConn.Create(&database.Friendship{
		Request:   false,
		Account:   accountId,
		Target:    target.String(),
		CreatedAt: time.Now().UnixMilli(),
	}).Error; err != nil {
		return fmt.Errorf("couldn't create friendship: %s", err)
	}

	// Notify the target to send them a notification
	if err := service.Instance.SendOne(standards.LiphiumAddress(target.String()).String(), neogate.Event{
		Name: "fr_rq",
		Data: requests.Map{
			"account":      accountLPH,
			"name":         accountInfo.Username,
			"display_name": accountInfo.DisplayName,
			"pub":          accountInfo.PublicKey,
			"sg":           accountInfo.SignatureKey,
		},
	}); err != nil {
		return fmt.Errorf("couldn't send event to target: %s", err)
	}

	return nil
}
