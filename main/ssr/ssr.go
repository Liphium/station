package ssr

import (
	"strings"

	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Reponse types
const ResponseRender = "render"
const ResponseRedirect = "redirect"
const ResponseError = "error"

// Redirect to a differnet path
func RedirectResponse(c *fiber.Ctx, path string, token string) fiber.Map {
	return fiber.Map{
		"success":  true,
		"type":     ResponseRedirect,
		"redirect": path,
	}
}

// Return an error on the current path that will be shown client-side above the submit button
func ErrorResponse(c *fiber.Ctx, err localization.Translations) fiber.Map {
	return fiber.Map{
		"success": true,
		"type":    ResponseError,
		"error":   translate(c, err),
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
	locale := locale(c)
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

// Extract the locale from any request
func locale(c *fiber.Ctx) string {
	locale := c.Locals("locale")
	if locale == nil {
		locale = localization.DefaultLocale
	}
	return strings.ToLower(locale.(string))
}

// Translate any message on a request
func translate(c *fiber.Ctx, message localization.Translations) string {
	locale := c.Locals("locale")
	if locale == nil {
		locale = localization.DefaultLocale
	}
	msg, valid := message[strings.ToLower(locale.(string))]
	if !valid {
		msg = message[localization.DefaultLocale]
	}
	return msg
}
