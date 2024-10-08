package sso_routes

import (
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
)

// Route: /account/auth/sso/form (from SSR redirect without token)
func getSSOForm(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token string `json:"token"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check the token
	state, msg := checkToken(req.Token)
	if msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Get the provider from goth
	provider, err := goth.GetProvider(openIdName)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Generate the session
	session, err := provider.BeginAuth(state.State)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get the url
	url, err := session.GetAuthURL()
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Add the goth session to the state
	state.Session = session.Marshal()

	// Return the SSO check form
	return util.ReturnJSON(c, ssr.RenderResponse(c, ssr.Components{
		ssr.Text{
			Text:  localization.RegisterSSOTitle,
			Style: ssr.TextStyleHeadline,
		},
		ssr.Text{
			Text:  localization.RegisterSSODescription,
			Style: ssr.TextStyleDescription,
		},
		ssr.Button{
			Link:  true,
			Label: localization.RegisterSSOButton,
			Path:  url,
		},
		ssr.StatusFetcher{
			Label:     localization.RegisterSSOStatus,
			Frequency: 5, // Refresh every 5 seconds
			Path:      "/account/auth/sso/check",
		},
	}))
}
