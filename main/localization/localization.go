package localization

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
