package util

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
func ConnectionToken(account uuid.UUID, session string, node uint) (string, error) {

	// Create jwt token
	exp := time.Now().Add(time.Hour * 2)
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, ConnectionTokenClaims{
		Account:        account.String(),
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
func SessionInformationToken(account uuid.UUID, sessions []string) (string, error) {

	// Create jwt token
	exp := time.Now().Add(time.Hour * 2)
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"acc": account.String(),
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
func Token(session uuid.UUID, account uuid.UUID, lvl uint, exp time.Time) (string, error) {

	// Create jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"ses": session.String(),
		"e_u": exp.Unix(), // Expiration unix
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

/*
This is staying here for now (even if not used) because we may need it
again in the future and because I might not want to recode this at that
time, this can slumber here until it is needed again (maybe).

If this isn't needed for another year, I think it can be removed. Currently
is the 14th of October in 2024. In case this project is still around in a few
years from now, let's celebrate by removing this piece of code!!! yay

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
*/
