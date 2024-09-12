package ssr

import (
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Reponse types
const ResponseRender = "render"
const ResponseRedirect = "redirect"
const ResponseSuggest = "suggest"
const ResponseSuccess = "success"

// Tell the client that the process is finished
func SuccessResponse(data interface{}) fiber.Map {
	return fiber.Map{
		"success": true,
		"type":    ResponseSuccess,
		"data":    data,
	}
}

// Redirect to a differnet path
func RedirectResponse(path string, token string) fiber.Map {
	return fiber.Map{
		"success":  true,
		"type":     ResponseRedirect,
		"redirect": path,
		"token":    token,
	}
}

// Response to suggest going to a different screen (for example register when email not found)
func SuggestResponse(c *fiber.Ctx, message localization.Translations, button Button) fiber.Map {
	locale := localization.Locale(c)
	return fiber.Map{
		"success": true,
		"type":    ResponseSuggest,
		"message": localization.TranslateLocale(locale, message),
		"button":  button.render(locale),
	}
}

// Render a new UI
func RenderResponse(c *fiber.Ctx, components Components) fiber.Map {
	// Catch an error if there is one
	defer func() {
		if err := recover(); err != nil {
			util.Log.Println("rendering error:", err)
		}
	}()

	// Render all the components
	locale := localization.Locale(c)
	util.Log.Println("using locale:", locale)
	compMap := make([]fiber.Map, len(components))
	for i, comp := range components {
		compMap[i] = comp.render(locale)
	}

	// Return a response with the rendered components
	return fiber.Map{
		"success": true,
		"type":    ResponseRender,
		"render":  compMap,
	}
}
