package auth

import (
	"node-backend/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateLoginTokenWithStep(id string, device string, step uint) (string, error) {

	// Create jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"s":   step,
		"e_u": time.Now().Add(time.Minute * 5).Unix(), // Expiration unix
		"acc": id,
		"d":   device,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tk.SignedString([]byte(util.JWT_SECRET))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetLoginDataFromToken(c *fiber.Ctx) (id string, device string, step uint) {

	// Get token from header
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	return claims["acc"].(string), claims["d"].(string), uint(claims["s"].(float64))
}
