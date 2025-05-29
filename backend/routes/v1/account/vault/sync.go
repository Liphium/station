package vault

import (
	"log"
	"sync"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/vault/sync
func syncVault(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Tags map[string]int64 `json:"tags"` // Tag -> Version
	}
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get account id
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}

	// Pull all of the entries and deletions in parallel
	wg := &sync.WaitGroup{}
	entryMap := &sync.Map{} // Tag -> []database.VaultEntry
	for tag, version := range req.Tags {
		wg.Add(1)

		// Start a goroutine for getting all the vault entries
		go func() {

			// Pull the entries for the tag with a version higher than currently
			var entries []database.VaultEntry
			if err := database.DBConn.Where("tag = ? AND account = ? AND version > ?", tag, accId, version).Find(&entries).Error; err != nil {
				log.Println("error while pulling entries:", err)
			} else {

				// Add all the entries to the map
				entryMap.Store(tag, entries)
			}

			wg.Done()
		}()
	}

	// Wait for all the data to arrive
	wg.Wait()

	// Collect all the results together
	allEntries := []database.VaultEntry{}
	for tag := range req.Tags {
		if entries, ok := entryMap.Load(tag); ok {
			allEntries = append(allEntries, entries.([]database.VaultEntry)...)
		} else {
			return integration.FailedRequest(c, localization.ErrorServer, nil)
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"entries": allEntries,
	})
}
