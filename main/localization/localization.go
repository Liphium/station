package localization

const DefaultLocale = "en_US"

// Predefined locales
var englishUS = "en_US"

// var german = "de_DE"

type Translations map[string]string

func None() Translations {
	return Translations{
		englishUS: "",
	}
}
