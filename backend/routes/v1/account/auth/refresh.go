package auth_routes

import (
	"errors"
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/requests"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type refreshRequest struct {
	Session string `json:"session"`
	Token   string `json:"token"`
}

// Route: /auth/refresh
func refreshSession(c *fiber.Ctx) error {

	// Parse request
	var req refreshRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Parse the session id
	id, err := uuid.Parse(req.Session)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Check if session is valid
	var session database.Session
	if !requests.GetSession(id, &session) {
		return util.ReturnJSON(c, fiber.Map{
			"success": false,
			"valid":   false,
		})
	}

	// Check if the session token matches the request
	if session.Token != req.Token {
		return util.ReturnJSON(c, fiber.Map{
			"success": false,
			"valid":   false,
		})
	}

	// Check if the session is verified
	if !session.Verified {
		var request database.KeyRequest = database.KeyRequest{
			Payload: "",
		}
		if err := database.DBConn.Where("session = ?", session.ID).Take(&request).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

			return util.ReturnJSON(c, fiber.Map{
				"success":  false,
				"verified": false,
			})
		}

		// Check if the key request has been accepted
		if request.Payload == "" {
			return util.ReturnJSON(c, fiber.Map{
				"success":  false,
				"verified": false,
			})
		}

		// Update the session to verified in case it has
		session.Verified = true // will be updated in the database anyway (below)
	}

	// Refresh session
	session.LastUsage = time.Now()
	if err := database.DBConn.Save(&session).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Create new token
	jwtToken, err := util.Token(session.ID, session.Account, session.PermissionLevel, time.Now().Add(time.Hour*24*3))

	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success":       true,
		"token":         jwtToken,
		"refresh_token": session.Token,
	})
}
