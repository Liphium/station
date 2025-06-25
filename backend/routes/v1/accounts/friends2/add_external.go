package friends2_routes

import (
	"fmt"

	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util/requests"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var NodeProtocol = "http://"

// Route: /accounts/friends/add_external
func addFriendFromExternal(c *fiber.Ctx) error {
	return nil
}

// Create a friend request from a target account id for an account. “accountLPH“ should be an address.
func createFriendRequest(accountLPH string, target uuid.UUID) error {
	accountId, origin, valid := standards.SplitLiphiumAddress(accountLPH)
	if !valid {
		return fmt.Errorf("invalid address: %s", accountLPH)
	}

	// Verify account in case not decentralized
	if origin != standards.CurrentTown() {
		res, err := requests.PostRequest(origin, "/accounts/get", requests.Map{
			"id": accountId,
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
	}

	return nil
}
