package auth

import (
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Request to sen
type registerFinishRequest struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Tag      string `json:"tag"`
	Password string `json:"password"`
}

// Route: /auth/register/finish, Finish the registration process
func registerFinish(c *fiber.Ctx) error {

	// Parse body to register request
	var req registerFinishRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Parse token and check token for validity
	tk, err := jwt.ParseWithClaims(req.Token, &registerClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(util.JWT_SECRET), nil
	})
	if err != nil {
		return util.InvalidRequest(c)
	}
	if !tk.Valid {
		return util.InvalidRequest(c)
	}

	// Check if email matches
	claims, ok := tk.Claims.(*registerClaims)
	if !ok {
		return util.InvalidRequest(c)
	}
	if claims.Step != 2 {
		return util.InvalidRequest(c)
	}
	if claims.ExpiredUnixSec < time.Now().Unix() {
		return util.InvalidRequest(c)
	}

	// Check if the email is already registered (needed here because the same email could technically be used (with another invite) to generate a
	// new JWT and send a new verification)
	if database.DBConn.Where("email = ?", claims.Email).Take(&account.Account{}).RowsAffected > 0 {
		return util.FailedRequest(c, "email.registered", nil)
	}

	// Check username and tag
	valid, message := standards.CheckUsernameAndTag(req.Username, req.Tag)
	if !valid {
		return util.FailedRequest(c, message, nil)
	}

	// Create account
	var acc account.Account = account.Account{
		Email:    claims.Email,
		Username: req.Username,
		RankID:   1, // Default rank
	}
	err = database.DBConn.Create(&acc).Error
	if err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// Create password authentication
	hash, err := auth.HashPassword(req.Password, acc.ID)
	if err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}
	err = database.DBConn.Create(&account.Authentication{
		ID:      auth.GenerateToken(8),
		Account: acc.ID,
		Type:    account.TypePassword,
		Secret:  hash,
	}).Error
	if err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}
	return util.SuccessfulRequest(c)
}
