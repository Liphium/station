package util

import (
	"os"
	"runtime/debug"

	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Auth errors
const EmailInvalid = "email.invalid"
const CodeInvalid = "code.invalid"
const PasswordInvalid = "password.incorrect"
const EmailRegistered = "email.registered" // When it is already registered
const UsernameInvalid = "username.invalid"
const UsernameTaken = "username.taken"
const TagInvalid = "tag.invalid"
const InviteInvalid = "invite.invalid"

func DebugRouteError(c *fiber.Ctx, msg localization.Translations) {
	if Testing {
		Log.Println(c.Route().Path+":", msg)
	}
}

func SuccessfulRequest(c *fiber.Ctx) error {
	return ReturnJSON(c, fiber.Map{
		"success": true,
	})
}

func FailedRequest(c *fiber.Ctx, message localization.Translations, err error) error {

	if LogErrors && err != nil {
		Log.Println(c.Route().Path+":", err)
		if os.Getenv("SHOW_STACK") == "1" {
			debug.PrintStack()
		}
	}

	return ReturnJSON(c, fiber.Map{
		"success": false,
		"error":   Translate(c, message),
	})
}

// Translate any message on a request
func Translate(c *fiber.Ctx, message localization.Translations) string {
	locale := c.Locals("locale").(string)
	if locale == "" {
		locale = localization.DefaultLocale
	}
	msg, _ := message[locale]
	return msg
}

func InvalidRequest(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusBadRequest)
}
