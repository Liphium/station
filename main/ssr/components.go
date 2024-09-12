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
	Name        string                    // Name in the return json
}

func (i Input) render(locale string) fiber.Map {
	return fiber.Map{
		"type":        "input",
		"placeholder": localization.TranslateLocale(locale, i.Placeholder),
		"hidden":      i.Hidden,
		"name":        i.Name,
	}
}

// The submit button, when this is clicked it's over
type Button struct {
	Label localization.Translations `json:"label,omitempty"` // Label on the button
	Path  string                    `json:"path"`            // The path the request goes to
}

func (b Button) render(locale string) fiber.Map {
	return fiber.Map{
		"type":  "button",
		"label": localization.TranslateLocale(locale, b.Label),
		"path":  b.Path,
	}
}

type Link struct {
	Label localization.Translations `json:"label,omitempty"` // Label on the button
}

func (l Link) render(locale string) fiber.Map {
	return fiber.Map{
		"type":  "link",
		"label": localization.TranslateLocale(locale, l.Label),
	}
}
