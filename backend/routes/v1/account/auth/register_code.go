package auth

import (
	"time"

	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Request sent to verify the email verification code
type registerCodeRequest struct {
	Token string `json:"token"`
	Code  string `json:"code"`
}

// Claims used for the third step in registration (generated by the second)
type registerClaims struct {
	Step           int    `json:"s"`
	Email          string `json:"e"`
	ExpiredUnixSec int64  `json:"e_u"`

	jwt.RegisteredClaims
}

// Route: /auth/register/code, Second step to registration, email code verification
func registerCode(c *fiber.Ctx) error {

	// Parse body to register request
	var req registerCodeRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Parse token and check token for validity
	tk, err := jwt.ParseWithClaims(req.Token, &registerEmailTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(util.JWT_SECRET), nil
	})
	if err != nil {
		return util.InvalidRequest(c)
	}
	if !tk.Valid {
		return util.InvalidRequest(c)
	}

	// Get claims of token
	claims, ok := tk.Claims.(*registerEmailTokenClaims)
	if !ok {
		return util.InvalidRequest(c)
	}
	if claims.Step != 1 {
		return util.InvalidRequest(c)
	}
	if claims.ExpiredUnixSec < time.Now().Unix() {
		return util.InvalidRequest(c)
	}

	// Unpack hidden value
	value, err := util.ReadHiddenJWTValue(c, claims.Code)
	if err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// Check code
	if req.Code != string(value) {
		return util.FailedRequest(c, util.CodeInvalid, err)
	}

	// Generate new token for the final step
	exp := time.Now().Add(time.Hour * 2)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS512, registerClaims{
		Step:           2,
		Email:          claims.Email,
		ExpiredUnixSec: exp.Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := newToken.SignedString([]byte(util.JWT_SECRET))
	if err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"token":   tokenString,
	})
}
