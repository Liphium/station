package vault

import "github.com/gofiber/fiber/v2"

// Configuration
const MaximumEntries = 1000 // Maximum number of entries in the vault (per account)

func Authorized(router fiber.Router) {
	router.Post("/add", addEntry)
	router.Post("/update", updateVaultEntry)
	router.Post("/remove", removeEntry)
	router.Post("/list", listEntries)
}
