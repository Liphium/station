package auth

import (
	"errors"
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LoginRequest is the request body for the login request
type startLoginRequest struct {
	Email  string `json:"email"`
	Device string `json:"device"`
}

// startLogin starts the login process (Route: /auth/login/start)
func startLogin(c *fiber.Ctx) error {

	var req startLoginRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if user exists
	var acc account.Account
	if database.DBConn.Where("email = ?", req.Email).Take(&acc).Error != nil {
		return util.FailedRequest(c, "email.invalid", nil)
	}

	valid, err := checkSessions(acc.ID)
	if err != nil {
		return util.FailedRequest(c, err.Error(), nil)
	}

	if !valid {
		return util.FailedRequest(c, "too.many.sessions", nil)
	}

	// Generate token
	return runAuthStep(acc.ID, req.Device, account.StartStep, c)
}

type loginStepRequest struct {
	Type   uint   `json:"type"`
	Secret string `json:"secret"`
}

// loginStep runs the login step (Route: /auth/login/step)
func loginStep(c *fiber.Ctx) error {

	var req loginStepRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get data
	id, device, step, err := auth.GetLoginDataFromToken(c)
	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	var method account.Authentication
	if err := database.DBConn.Where("account = ? AND type = ?", id, req.Type).Take(&method).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Check the provided secret
	if !method.Verify(req.Type, req.Secret, id) {
		return util.FailedRequest(c, "invalid.method", nil)
	}

	return runAuthStep(id, device, step+1, c)
}

// Runs the next step in an authentication
func runAuthStep(id uuid.UUID, device string, step uint, c *fiber.Ctx) error {

	// Generate token
	tk, err := auth.GenerateLoginTokenWithStep(id, device, step)
	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Get authentication methods for step
	var availableMethods []uint
	for method, stepUsed := range account.Order {
		if stepUsed == step {
			availableMethods = append(availableMethods, method)
		}
	}

	var methods []uint
	query := database.DBConn.Model(&account.Authentication{}).Where("type IN ? AND account = ?", availableMethods, id).Select("type").Take(&methods)

	if query.Error == gorm.ErrRecordNotFound && step == account.StartStep {
		// TODO: SERIOUS SECURITY ISSUE WARNING HERE
		return util.FailedRequest(c, "no.methods", nil)
	}

	if query.Error == gorm.ErrRecordNotFound {

		var acc account.Account
		if err := database.DBConn.Where("id = ?", id).Preload("Rank").Take(&acc).Error; err != nil {
			return util.FailedRequest(c, "server.error", err)
		}

		// Create session
		tk := auth.GenerateToken(100)

		var createdSession account.Session = account.Session{
			ID:              auth.GenerateToken(8),
			Token:           tk,
			Verified:        false,
			Account:         acc.ID,
			PermissionLevel: acc.Rank.Level,
			Device:          device,
			LastConnection:  time.UnixMilli(0),
		}

		if err := database.DBConn.Create(&createdSession).Error; err != nil {
			util.FailedRequest(c, "server.error", err)
		}

		// Generate jwt token
		jwtToken, err := util.Token(createdSession.ID, acc.ID, acc.Rank.Level, time.Now().Add(time.Hour*24*1))

		if err != nil {
			return util.FailedRequest(c, "server.error", err)
		}

		return util.ReturnJSON(c, fiber.Map{
			"success":       true,
			"token":         jwtToken,
			"refresh_token": tk,
		})
	}

	if query.Error != nil {
		return util.FailedRequest(c, "server.error", nil)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"token":   tk,
		"methods": methods,
	})
}

// checkSessions checks if the user has too many sessions
func checkSessions(id uuid.UUID) (bool, error) {

	// Check if user has too many sessions
	var sessions int64
	if err := database.DBConn.Model(&account.Session{}).Where("account = ?", id).Count(&sessions).Error; err != nil {
		return false, errors.New("server.error")
	}

	if sessions > 5 {
		return false, errors.New("sessions.limit")
	}

	return true, nil
}
