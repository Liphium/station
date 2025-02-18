package vault

import (
	"log"
	"sync"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/vault/sync
func syncVault(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Tags map[string]int64 `json:"tags"` // Tag -> Version
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get account id
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Pull all of the entries and deletions in parallel
	wg := &sync.WaitGroup{}
	entryMap := &sync.Map{}    // Tag -> []database.VaultEntry
	deletionMap := &sync.Map{} // Tag -> []string
	for tag, version := range req.Tags {
		wg.Add(2)

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

		// Start a goroutine to pull all the deletions for that tag
		go func() {

			// Pull the entries for the tag with a version higher than currently
			var deletions []string
			if err := database.DBConn.Model(&database.VaultDeletion{}).Select("entry").Where("tag = ? AND account = ? AND version > ?", tag, accId, version).Find(&deletions).Error; err != nil {
				log.Println("error while pulling deletions:", err)
			} else {

				// Add all the deletions to the map
				deletionMap.Store(tag, deletions)
			}

			wg.Done()
		}()
	}

	// Wait for all the data to arrive
	wg.Wait()

	// Collect all the results together
	entryMapJS := map[string][]database.VaultEntry{}
	deletionMapJS := map[string][]string{}
	for tag := range req.Tags {
		if entries, ok := entryMap.Load(tag); ok {
			entryMapJS[tag] = entries.([]database.VaultEntry)
		} else {
			return util.FailedRequest(c, localization.ErrorServer, nil)
		}
		if deletions, ok := deletionMap.Load(tag); ok {
			deletionMapJS[tag] = deletions.([]string)
		} else {
			return util.FailedRequest(c, localization.ErrorServer, nil)
		}
	}

	return util.ReturnJSON(c, fiber.Map{
		"success":   true,
		"entries":   entryMapJS,
		"deletions": deletionMapJS,
	})
}
