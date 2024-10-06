package townhall_settings

import (
	"github.com/Liphium/station/backend/settings"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /townhall/settings/files
func decentralizationSettings(c *fiber.Ctx) error {
	locale := localization.Locale(c)
	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"settings": []fiber.Map{
			settings.DecentralizationEnabled.ToMap(locale),
			settings.DecentralizationAllowUnsafe.ToMap(locale),
		},
	})
}
