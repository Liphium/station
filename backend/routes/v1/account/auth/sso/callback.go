package sso_routes

import (
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
)

// Route: /account/auth/sso/callback?code=%code% (after SSO)
func callback(c *fiber.Ctx) error {

	// Get the query parameters from SSO
	token := c.Query("state") // Equal to the SSR token

	// Validate the token
	state, msg := checkToken(token)
	if msg != nil {
		return c.SendString(localization.TranslateLocale(localization.DefaultLocale, msg))
	}

	// Get the provider
	provider, err := goth.GetProvider(openIdName)
	if err != nil {
		return c.SendString(localization.TranslateLocale(state.Locale, localization.ErrorServer))
	}

	// Unmarshal the session
	session, err := provider.UnmarshalSession(state.Session)
	if err != nil {
		return c.SendString(localization.TranslateLocale(state.Locale, localization.ErrorServer))
	}

	// Authorize using code and state
	_, err = session.Authorize(provider, &Params{
		ctx: c,
	})
	if err != nil {
		return c.SendString(localization.TranslateLocale(state.Locale, localization.ErrorServer))
	}

	// Get the user data (email and stuff)
	user, err := provider.FetchUser(session)
	if err != nil {
		return c.SendString(localization.TranslateLocale(state.Locale, localization.ErrorServer))
	}

	// Check the email
	valid, normalizedEmail := standards.CheckEmail(user.Email)
	if !valid {
		return c.SendString(localization.TranslateLocale(state.Locale, localization.ErrorEmailNotFound))
	}

	// Set the email and user id in the state (so it can be used in the check endpoint)
	state.UserID = user.UserID
	state.Email = normalizedEmail
	state.Done = true

	return c.SendString(localization.TranslateLocale(state.Locale, localization.RegisterSSOComplete))
}
