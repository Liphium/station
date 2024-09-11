package ssr

import (
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type Component interface {
	render(string) fiber.Map // The function that translates all the things and turns the component into a proper json
}
type Components []Component

type Text struct {
	Text  localization.Translations // The text itself
	Style int                       // Style of the text (0 = headline, 1 = description)
}

func (t Text) render(locale string) fiber.Map {
	return fiber.Map{
		"text":  t.Text[locale],
		"style": t.Style,
	}
}

type Input struct {
	Label       localization.Translations // Label above the input
	Placeholder localization.Translations // Placeholder inside the input on the client
	Name        string                    // Name in the return json
}

func (i Input) render(locale string) fiber.Map {
	return fiber.Map{
		"label":       i.Label[locale],
		"placeholder": i.Placeholder[locale],
		"name":        i.Name,
	}
}

// The submit button, when this is clicked it's over
type Button struct {
	Label localization.Translations `json:"label,omitempty"` // Label on the button
	Path  string                    `json:"path"`            // The path the request goes to
}

func (b *Button) render(locale string) fiber.Map {
	return fiber.Map{
		"label": b.Label[locale],
		"path":  b.Path,
	}
}

type Link struct {
	Label localization.Translations `json:"label,omitempty"` // Label on the button
}

func (l Link) render(locale string) fiber.Map {
	return fiber.Map{
		"label": l.Label[locale],
	}
}
