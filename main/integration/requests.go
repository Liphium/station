package integration

import (
	"runtime/debug"

	"github.com/Liphium/station/main/localization"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func SuccessfulRequest(c *fiber.Ctx) error {
	return ReturnJSON(c, fiber.Map{
		"success": true,
	})
}

func FailedRequest(c *fiber.Ctx, message localization.Translations, err error) error {

	// Print error if it isn't nil
	if err != nil {
		Log.Println(c.Route().Name + " ERROR: " + message[localization.DefaultLocale] + ":" + err.Error())
		debug.PrintStack()
	}

	return ReturnJSON(c, fiber.Map{
		"success": false,
		"error":   Translate(c, message),
	})
}

func InvalidRequest(c *fiber.Ctx, message string) error {
	Log.Println(c.Route().Name + " request is invalid. msg: " + message)
	debug.PrintStack()
	return c.SendStatus(fiber.StatusBadRequest)
}

// Parse encrypted json
func BodyParser(c *fiber.Ctx, data interface{}) error {
	return sonic.Unmarshal(c.Locals("body").([]byte), data)
}

// Translate any message on a request
func Translate(c *fiber.Ctx, message localization.Translations) string {
	locale := c.Locals("locale")
	if locale == nil {
		locale = localization.DefaultLocale
	}
	msg := message[locale.(string)]
	return msg
}

// Return encrypted json
func ReturnJSON(c *fiber.Ctx, data interface{}) error {

	encoded, err := sonic.Marshal(data)
	if err != nil {
		return FailedRequest(c, localization.ErrorServer, err)
	}

	if c.Locals("key") == nil {
		return c.Send(encoded)
	}
	encrypted, err := EncryptAES(c.Locals("key").([]byte), encoded)
	if err != nil {
		return FailedRequest(c, localization.ErrorServer, err)
	}

	return c.Send(encrypted)
}
