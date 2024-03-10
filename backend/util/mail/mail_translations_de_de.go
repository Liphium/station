package mail

import (
	"node-backend/util"
	"os"
)

// English is the default for email translations yk
var deDeTranslations = map[string]func([]string) (string, []string){
	EmailVerification: func(args []string) (string, []string) {
		appName := os.Getenv(util.EnvAppName)
		return appName + " Verifizierungscode für deine Email", []string{
			"Hallo von " + appName + "!",
			" ",
			"Wir schreiben dir, um deine Email zu verifizieren und so, weil du einen Account gemacht hast! Wenn du das nicht warst, haben wir hier vielleicht einen kleines Problem vorliegen. Bitte geh dann zu unserer Webseite und melde dieses Problem.",
			"Naja, hier ist dein Verifizierungscode: " + args[0],
			" ",
			"Danke fürs Registrieren,",
			"dein " + appName + " Team.",
		}
	},
}
