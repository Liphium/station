package register_routes

import (
	"time"

	"github.com/Liphium/station/backend/database"
	login_routes "github.com/Liphium/station/backend/routes/v1/account/auth/login"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/register/password
func checkPassword(c *fiber.Ctx) error {

	// Parse the basic request (it could be SSO and not contain password and stuff)
	var basicReq struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&basicReq); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Check the token
	state, msg := validateToken(basicReq.Token, 6)
	if msg != nil {
		return integration.FailedRequest(c, msg, nil)
	}

	password := ""
	if !state.SSO {
		// If SSO isn't enabled, check the password and stuff
		var req struct {
			Token           string `json:"token"`
			Password        string `json:"password"`
			ConfirmPassword string `json:"confirm_password"`
		}
		if err := c.BodyParser(&req); err != nil {
			return integration.InvalidRequest(c, "invalid request")
		}

		// Rate limit the amount of requests
		if !ratelimitHandler(state, 2, time.Second*10) {
			return integration.FailedRequest(c, localization.ErrorAuthRatelimit, nil)
		}

		// Check the requirements
		if len(req.Password) < 8 {
			return integration.FailedRequest(c, localization.ErrorPasswordInvalid(8), nil)
		}
		if req.Password != req.ConfirmPassword {
			return integration.FailedRequest(c, localization.ErrorPasswordsDontMatch, nil)
		}
		password = req.Password
	}

	// Re-check all of the data
	if msg := standards.CheckUsername(state.Username); msg != nil {
		return integration.FailedRequest(c, localization.ErrorRegistrationFailed(msg), nil)
	}
	if msg := standards.CheckDisplayName(state.DisplayName); msg != nil {
		return integration.FailedRequest(c, localization.ErrorRegistrationFailed(msg), nil)
	}
	if err := database.DBConn.Where("email = ?", state.Email).Take(&database.Account{}).Error; err == nil {
		return integration.FailedRequest(c, localization.ErrorRegistrationFailed(localization.ErrorEmailAlreadyInUse), nil)
	}

	// Check if there are other accounts
	var count int64
	if err := database.DBConn.Model(&database.Account{}).Count(&count).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get the rank the user should have
	var rankToGive database.Rank
	if count == 0 {
		if err := database.DBConn.Where("name = ?", "Admin").Take(&rankToGive).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
	} else {
		if err := database.DBConn.Where("name = ?", "Default").Take(&rankToGive).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
	}

	// Redeem the invite code (not when using SSO)
	if !state.SSO {

		// Make sure the invite exists in the first place
		if err := database.DBConn.Where("id = ?", state.Invite).Take(&database.Invite{}).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorRegistrationFailed(localization.ErrorInviteInvalid), err)
		}

		if err := database.DBConn.Where("id = ?", state.Invite).Delete(&database.Invite{}).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
	}

	// Create an account
	acc := database.Account{
		Email:       state.Email,
		Username:    state.Username,
		DisplayName: state.DisplayName,
		RankID:      rankToGive.ID,
	}
	if err := database.DBConn.Create(&acc).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Create authentication
	if state.SSO {

		// Create SSO-based authentication when the server is using SSO
		if err := database.DBConn.Create(&database.Authentication{
			Account: acc.ID,
			Type:    database.AuthTypeSSO,
			Secret:  state.UserID,
		}).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
	} else {

		// Create password-based authentication if the server isn't using SSO
		hash, err := auth.HashPassword(password, acc.ID)
		if err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
		if err := database.DBConn.Create(&database.Authentication{
			Account: acc.ID,
			Type:    database.AuthTypePassword,
			Secret:  hash,
		}).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
	}

	// Create the tokens
	token, refreshToken, err := login_routes.CreateSession(acc.ID, rankToGive.Level)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.JSON(ssr.SuccessResponse(fiber.Map{
		"token":         token,
		"refresh_token": refreshToken,
	}))
}
