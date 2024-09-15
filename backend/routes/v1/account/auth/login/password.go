package login_routes

import (
	"errors"
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/kv"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Route: /account/auth/login/password (SSR)
func checkPassword(c *fiber.Ctx) error {

	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if there is a login token like this
	obj, valid := kv.Get(loginTokenPrefix + req.Token)
	if !valid {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	// Make sure the state is valid
	state := obj.(*LoginState)
	if state.LoginStep != 1 {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	// Handle the rate limit
	if !ratelimitHandler(state) {
		return util.FailedRequest(c, localization.ErrorAuthRatelimit, nil)
	}

	// Get the password stored in the database
	var credential account.Authentication
	if err := database.DBConn.Where("account = ? AND type = ?", state.Account, account.TypePassword).Take(&credential).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Check the password
	match, err := auth.ComparePasswordAndHash(req.Password, state.Account, credential.Secret)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}
	if !match {
		return util.FailedRequest(c, localization.ErrorPasswordIncorrect, nil)
	}

	// Remove the login token from the kv
	kv.Delete(loginTokenPrefix + req.Token)

	// Count the amount of sessions
	var sessionCount int64 = 0
	if err := database.DBConn.Model(&account.Session{}).Where("account = ?", state.Account).Count(&sessionCount).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Create session
	tk := auth.GenerateToken(100)
	var createdSession account.Session = account.Session{
		Token:           tk,
		Verified:        sessionCount == 0,
		Account:         state.Account,
		PermissionLevel: state.PermissionLevel,
		Device:          "tbd",
		LastConnection:  time.UnixMilli(0),
	}

	// Create the session in a safe way
	if err = database.DBConn.Create(&createdSession).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Generate jwt token for the session
	jwtToken, err := util.Token(createdSession.ID, state.Account, state.PermissionLevel, time.Now().Add(time.Hour*24*1))
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return refresh and normal token
	return util.ReturnJSON(c, ssr.SuccessResponse(fiber.Map{
		"token":         jwtToken,
		"refresh_token": tk,
	}))
}
