package routes

import (
	"strings"

	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
	"github.com/gofiber/fiber/v2"
)

// Route: /create
func createSpace(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token string `json:"token"`
	}
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "request was invalid")
	}

	// Validate the token
	claims, valid := caching.SSInstance.CheckToken(req.Token, caching.SSNode)
	if !valid {
		return integration.InvalidRequest(c, "token was invalid")
	}

	// Make sure it's not an only join token
	if strings.HasPrefix(claims.Extra, "oj-") {
		return integration.InvalidRequest(c, "no permission")
	}

	// Create a new Space
	spaceId := util.GenerateToken(32)
	for {
		if _, valid := caching.GetRoom(spaceId); !valid {
			break
		}
		spaceId = util.GenerateToken(32)
	}
	caching.CreateRoom(spaceId)

	// Create a connection token
	clientId, token, err := newConnectionToken(c, spaceId)
	if err != nil {
		return err
	}

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"client":  clientId,
		"space":   spaceId,
		"token":   token,
	})
}
