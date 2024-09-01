package remote_action_routes

import (
	"errors"
	"sync"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

// The map storing all the tokens and their corresponding nodes (Token ID -> Node address)
var tokenMap *sync.Map = &sync.Map{}

// Action: negotiate
func handleNegotiation(c *fiber.Ctx, action struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Node  string `json:"node"`
}) error {

	// Validate the token
	_, err := caching.ValidateToken(action.ID, action.Token)
	if err != nil {

		// If the token is invalid, return invalid credentials
		if errors.Is(err, caching.ErrInvalidToken) {
			return integration.FailedRequest(c, localization.InvalidCredentials, nil)
		}

		// Otherwise, return a normal server error
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get the sender
	if action.Node == "" || action.Node == util.OwnPath {
		return integration.InvalidRequest(c, "invalid sender")
	}

	// Save the node to the map
	tokenMap.Store(action.ID, action.Node)

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"node":    util.OwnPath,
	})
}
