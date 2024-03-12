package util

import (
	"time"

	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RemoteId(lvl uint) (string, error) {

	// Create jwt token
	exp := time.Now().Add(time.Hour * 2)
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"e_u": exp.Unix(), // Expiration unix
		"lvl": lvl,
		"rid": true, // tell the backend that it's a remote id
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tk.SignedString([]byte(integration.JwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// IsExpired checks if the token is expired
func IsExpired(c *fiber.Ctx) bool {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	num := claims["e_u"].(float64)
	exp := int64(num)

	return time.Now().Unix() > exp
}

// Permission checks if the user has the required permission level
func Permission(c *fiber.Ctx, perm int16) bool {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	lvl := int16(claims["lvl"].(float64))

	return lvl >= perm
}
