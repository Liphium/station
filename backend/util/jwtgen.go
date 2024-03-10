package util

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Connection token struct
type ConnectionTokenClaims struct {
	Account        string `json:"acc"`  // Account id of the connecting client
	ExpiredUnixSec int64  `json:"e_u"`  // Expiration time in unix seconds
	Session        string `json:"ses"`  // Session id of the connecting client
	Node           string `json:"node"` // Node id of the node the client is connecting to

	jwt.RegisteredClaims
}

// Generate a connection token for a node
func ConnectionToken(account string, session string, node uint) (string, error) {

	// Create jwt token
	exp := time.Now().Add(time.Hour * 2)
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, ConnectionTokenClaims{
		Account:        account,
		ExpiredUnixSec: exp.Unix(),
		Session:        session,
		Node:           fmt.Sprintf("%d", node),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tk.SignedString([]byte(JWT_SECRET))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Create a token with current session information (some nodes may require this)
func SessionInformationToken(account string, sessions []string) (string, error) {

	// Create jwt token
	exp := time.Now().Add(time.Hour * 2)
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"acc": account,
		"e_u": exp.Unix(), // Expiration unix
		"se":  sessions,   // Session list (for the node)
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tk.SignedString([]byte(JWT_SECRET))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Generate a normal authenticated token
func Token(session string, account string, lvl uint, exp time.Time) (string, error) {

	// Create jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"ses": session,
		"e_u": exp.Unix(), // Expiration unix
		"acc": account,
		"lvl": lvl,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tk.SignedString([]byte(JWT_SECRET))

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

	// Check if it is actually a node token (shouldn't be usable here)
	_, valid := claims["node"]
	if valid {
		return false
	}

	return time.Now().Unix() > exp
}

// Permission checks if the user has the required permission level
func Permission(c *fiber.Ctx, perm string) bool {

	// Check if there is a JWT token
	if c.Locals("user") == nil || reflect.TypeOf(c.Locals("user")).String() != "*jwt.Token" {
		return false
	}

	// Get the permission from the map
	permission, valid := Permissions[perm]
	if !valid {
		return false
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	lvl := int16(claims["lvl"].(float64))

	return lvl >= permission
}

func GetSession(c *fiber.Ctx) string {
	if c.Locals("user") == nil || reflect.TypeOf(c.Locals("user")).String() != "*jwt.Token" {
		return ""
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["ses"].(string)
}

func GetAcc(c *fiber.Ctx) string {
	if c.Locals("user") == nil || reflect.TypeOf(c.Locals("user")).String() != "*jwt.Token" {
		return ""
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["acc"].(string)
}

// Generate a JWT value that the client can't read (can't be really long because of RSA encryption)
func MakeHiddenJWTValue(c *fiber.Ctx, value []byte) (string, error) {
	pub := c.Locals(LocalsServerPub).(*rsa.PublicKey)

	// Encrypt with RSA
	encrypted, err := EncryptRSA(pub, value)
	if err != nil {
		return "", err
	}

	// Encode for use in JSON
	encoded := base64.StdEncoding.EncodeToString(encrypted)
	return encoded, nil
}

// Read a "hidden" JWT value encrypted by the server (referred to as a "hidden value")
func ReadHiddenJWTValue(c *fiber.Ctx, encoded string) ([]byte, error) {
	priv := c.Locals(LocalsServerPriv).(*rsa.PrivateKey)

	// Decode to bytes
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	// Decrypt with RSA
	decrypted, err := DecryptRSA(priv, decoded)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
