package auth

import (
	"time"

	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateLoginTokenWithStep(id uuid.UUID, device string, step uint) (string, error) {

	// Create jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"s":   step,
		"e_u": time.Now().Add(time.Minute * 5).Unix(), // Expiration unix
		"acc": id.String(),
		"d":   device,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tk.SignedString([]byte(util.JWT_SECRET))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetLoginDataFromToken(c *fiber.Ctx) (id uuid.UUID, device string, step uint, parseErr error) {

	// Get token from header
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	// Parse uuid
	id, err := uuid.Parse(claims["acc"].(string))
	if err != nil {
		return uuid.UUID{}, "", 0, err
	}

	return id, claims["d"].(string), uint(claims["s"].(float64)), nil
}
