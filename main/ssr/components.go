package ssr

import (
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type Component interface {
	render(string) fiber.Map // The function that translates all the things and turns the component into a proper json
}
type Components []Component

const TextStyleHeadline = 0    // Style for the headline above the form
const TextStyleSeperator = 1   // Style for a seperator between sections
const TextStyleDescription = 2 // Style for a description or normal text element

type Text struct {
	Text  localization.Translations // The text itself
	Style int                       // Style of the text (0 = headline, 1 = description)
}

func (t Text) render(locale string) fiber.Map {
	return fiber.Map{
		"type":  "text",
		"text":  localization.TranslateLocale(locale, t.Text),
		"style": t.Style,
	}
}

type Input struct {
	Placeholder localization.Translations // Placeholder inside the input on the client
	Hidden      bool                      // If the characters inside the input should be hidden
	UTF8        bool                      // If the string should be encoded using UTF8
	Value       string                    // A pre-filled value already in the input
	Name        string                    // Name in the return json
}

func (i Input) render(locale string) fiber.Map {
	return fiber.Map{
		"type":        "input",
		"placeholder": localization.TranslateLocale(locale, i.Placeholder),
		"hidden":      i.Hidden,
		"utf8":        i.UTF8,
		"value":       i.Value,
		"name":        i.Name,
	}
}

// The submit button, when this is clicked it's over
type SubmitButton struct {
	Label localization.Translations // Label on the button
	Path  string                    // The path the request goes to
}

func (b SubmitButton) render(locale string) fiber.Map {
	return fiber.Map{
		"type":  "submit",
		"label": localization.TranslateLocale(locale, b.Label),
		"path":  b.Path,
	}
}

// A regular button, when this is clicked an action can be completed (popup, error, etc.)
type Button struct {
	Label localization.Translations // Label on the button
	Link  bool                      // If the button is actually a link (to a website)
	Path  string                    // The path the request goes to
}

func (b Button) render(locale string) fiber.Map {
	return fiber.Map{
		"type":  "submit",
		"link":  b.Link,
		"label": localization.TranslateLocale(locale, b.Label),
		"path":  b.Path,
	}
}
