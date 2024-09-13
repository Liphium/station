package sso_routes

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
)

// Route: /account/auth/sso/redirect
func beginAuth(c *fiber.Ctx) error {

	// Get the provider from goth
	provider, err := goth.GetProvider(openIdName)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Generate the session
	session, err := provider.BeginAuth(getState(c))
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get the url
	url, err := session.GetAuthURL()
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

// This function was copied:
func getState(c *fiber.Ctx) string {
	// Check if there is already a state param
	state := c.Query("state")
	if len(state) > 0 {
		return state
	}

	// If a state query param is not passed in, generate a random
	// base64-encoded nonce so that the state on the auth URL
	// is unguessable, preventing CSRF attacks, as described in
	//
	// https://auth0.com/docs/protocols/oauth2/oauth-state#keep-reading
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic("gothic: source of randomness unavailable: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}
