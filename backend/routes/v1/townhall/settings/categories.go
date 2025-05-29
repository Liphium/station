package townhall_settings

import (
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type category struct {
	Name localization.Translations
	ID   string
}

var categories = []category{
	{
		Name: localization.SettingCategoryFiles,
		ID:   "files",
	},
	{
		Name: localization.SettingCategoryChat,
		ID:   "decentralization",
	},
}

// Route: /townhall/settings/categories
func getCategories(c *fiber.Ctx) error {

	// Convert the categories to a list with translation
	locale := localization.Locale(c)
	translated := make([]fiber.Map, len(categories))
	for i, c := range categories {
		translated[i] = fiber.Map{
			"name": localization.TranslateLocale(locale, c.Name),
			"id":   c.ID,
		}
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"categories": translated,
	})
}
