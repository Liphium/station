package util

import (
	"time"

	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type timestampClaims struct {
	Creation int64 `json:"c"`
	jwt.MapClaims
}

func TimestampToken(time int64) (string, error) {

	// Create jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, timestampClaims{
		Creation: time,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tk.SignedString([]byte(integration.JwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyTimestampToken(timestampToken string) (int64, bool) {
	token, err := jwt.ParseWithClaims(timestampToken, &timestampClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(integration.JwtSecret), nil
	}, jwt.WithLeeway(5*time.Minute))

	if err != nil {
		return 0, false
	}

	if claims, ok := token.Claims.(*timestampClaims); ok && token.Valid {
		return claims.Creation, true
	}

	return 0, false
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
