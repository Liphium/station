package recovery_routes

import "github.com/gofiber/fiber/v2"

// Flow of adding a new recovery key:
//
// 1. Client generates a random key and encrypts keys with it.
//
// 2. Calls /generate with the encrypted keys and gives the user a token in the format:
// SERVER_TOKEN-ENCRYPTION_KEY
//
// This ensures that the server never gets the key required to decrypt, but can still verify
// that the user actually used a recovery token valid to the server (for session verification).

// Configuration
const MaxRecoveryTokens = 5

func Authorized(router fiber.Router) {
	router.Post("/generate", generateRecoveryToken)
	router.Post("/list", listRecoveryTokens)
	router.Post("/delete", deleteRecoveryToken)
}
