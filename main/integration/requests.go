package integration

import (
	"runtime/debug"
	"strings"

	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

func SuccessfulRequest(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
	})
}

func FailedRequest(c *fiber.Ctx, message localization.Translations, err error) error {

	// Print error if it isn't nil
	if err != nil {
		Log.Println(c.Route().Name + " ERROR: " + message[localization.DefaultLocale] + ":" + err.Error())
		debug.PrintStack()
	}

	return c.JSON(fiber.Map{
		"success": false,
		"error":   Translate(c, message),
	})
}

func InvalidRequest(c *fiber.Ctx, message string) error {
	Log.Println(c.Route().Name + " request is invalid. msg: " + message)
	debug.PrintStack()
	return c.SendStatus(fiber.StatusBadRequest)
}

// Translate any message on a request
func Translate(c *fiber.Ctx, message localization.Translations) string {
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
