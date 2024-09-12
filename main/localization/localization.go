package localization

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

const DefaultLocale = "en_us"

// Predefined locales
var englishUS = "en_us"

// var german = "de_DE"

type Translations map[string]string

func None() Translations {
	return Translations{
		englishUS: "",
	}
}

// Extract the locale from any request
func Locale(c *fiber.Ctx) string {
	locale := c.Locals("locale")
	if locale == nil {
		locale = DefaultLocale
	}
	return strings.ToLower(locale.(string))
}

// Translate any message using a locale
func TranslateLocale(locale string, message Translations) string {
	msg, valid := message[locale]
	if !valid {
		msg = message[DefaultLocale]
	}
	return msg
}

// Translate any message on a request
func Translate(c *fiber.Ctx, message Translations) string {
	return TranslateLocale(Locale(c), message)
}
