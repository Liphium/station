package message_handlers

import (
	"log"
	"time"

	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func SetupHandler() {
	// Handlers for sending messages
	pipeshandler.CreateHandlerFor(caching.SSInstance, "msg_timestamp", generateTimestampToken)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "msg_send", sendMessage)
}

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
	})
	if err != nil {
		log.Println(timestampToken, err)
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
