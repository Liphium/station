package register_routes

import (
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/register/password
func checkPassword(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token           string `json:"token"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check the token
	state, msg := validateToken(req.Token, 5)
	if msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Rate limit the amount of requests
	if !ratelimitHandler(state, 2, time.Second*10) {
		return util.FailedRequest(c, localization.ErrorAuthRatelimit, nil)
	}

	// Check the requirements
	if len(req.Password) < 8 {
		return util.FailedRequest(c, localization.ErrorPasswordInvalid(8), nil)
	}
	if req.Password != req.ConfirmPassword {
		return util.FailedRequest(c, localization.ErrorPasswordsDontMatch, nil)
	}

	// Re-check all of the data
	if msg := standards.CheckUsername(state.Username); msg != nil {
		return util.FailedRequest(c, localization.ErrorRegistrationFailed(msg), nil)
	}
	if msg := standards.CheckDisplayName(state.DisplayName); msg != nil {
		return util.FailedRequest(c, localization.ErrorRegistrationFailed(msg), nil)
	}
	if err := database.DBConn.Where("email = ?", state.Email).Take(&account.Account{}).Error; err == nil {
		return util.FailedRequest(c, localization.ErrorRegistrationFailed(localization.ErrorEmailAlreadyInUse), nil)
	}

	// Get the default rank
	var defaultRank account.Rank
	if err := database.DBConn.Where("id = ?", 1).Take(&defaultRank).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Redeem the invite code
	if err := database.DBConn.Where("id = ?", state.Invite).Delete(&account.Invite{}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Create an account
	acc := account.Account{
		Email:       state.Email,
		Username:    state.Username,
		DisplayName: state.DisplayName,
		RankID:      1, // Default rank
	}
	if err := database.DBConn.Create(&acc).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Create password authentication
	hash, err := auth.HashPassword(req.Password, acc.ID)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}
	if err := database.DBConn.Create(&account.Authentication{
		Account: acc.ID,
		Type:    account.TypePassword,
		Secret:  hash,
	}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Create session
	tk := auth.GenerateToken(100)
	var createdSession account.Session = account.Session{
		Token:           tk,
		Verified:        true,
		Account:         acc.ID,
		PermissionLevel: defaultRank.Level,
		Device:          "tbd",
		LastConnection:  time.UnixMilli(0),
	}

	// Create the session in a safe way
	if err = database.DBConn.Create(&createdSession).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Generate jwt token for the session
	jwtToken, err := util.Token(createdSession.ID, acc.ID, defaultRank.Level, time.Now().Add(time.Hour*24*1))
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, ssr.SuccessResponse(fiber.Map{
		"token":         jwtToken,
		"refresh_token": createdSession.Token,
	}))
}
