package remote_action_routes

import (
	"sync"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
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
		return integration.FailedRequest(c, localization.ErrorInvalidCredentials, nil)
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
