package friends2_routes

import (
	"fmt"

	"github.com/Liphium/station/backend/standards"
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
	_, origin, valid := standards.SplitLiphiumAddress(accountLPH)
	if !valid {
		return fmt.Errorf("invalid address: %s", accountLPH)
	}

	// Verify account in case not decentralized
	if origin != standards.CurrentTown() {

	}

	return nil
}
