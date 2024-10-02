package sso_routes

import (
	"github.com/Liphium/station/backend/database"
	login_routes "github.com/Liphium/station/backend/routes/v1/account/auth/login"
	register_routes "github.com/Liphium/station/backend/routes/v1/account/auth/register"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/sso/check (from SSR status fetcher)
func checkSSO(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token string `json:"token"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Validate the token
	state, msg := checkToken(req.Token)
	if msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Also return a failed request if SSO hasn't been completed yet
	if !state.Done {
		msg = localization.ErrorSSONotCompleted
	}
	if msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Check if there is an account with that user id already
	var auth database.Authentication
	err := database.DBConn.Where("secret = ? AND type = ?", state.UserID, database.AuthTypeSSO).Take(&auth).Error
	if err == nil {

		// Get the account from the authentication method if one exists
		var acc database.Account
		if err := database.DBConn.Where("id = ?", auth.Account).Preload("Rank").Take(&acc).Error; err != nil {
			return util.FailedRequest(c, localization.ErrorServer, err)
		}

		// Generate a new session for this account
		token, refreshToken, err := login_routes.CreateSession(acc.ID, acc.Rank.Level)
		if err != nil {
			return util.FailedRequest(c, localization.ErrorServer, err)
		}

		return util.ReturnJSON(c, ssr.SuccessResponse(fiber.Map{
			"token":         token,
			"refresh_token": refreshToken,
		}))
	}

	// Create a new register token with SSO enabled to circumvent the password prompt
	token := register_routes.GenerateRegisterTokenForSSO(state.Email, state.UserID)
	return util.ReturnJSON(c, ssr.RedirectResponse("/account/auth/register/from_sso", token))
}
