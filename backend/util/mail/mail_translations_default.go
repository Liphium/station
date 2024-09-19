package mail

import (
	"os"

	"github.com/Liphium/station/backend/util"
)

// English is the default for email translations yk
var enUsTranslations = map[string]func([]string) (string, []string){
	EmailVerification: func(args []string) (string, []string) {
		appName := os.Getenv(util.EnvAppName)
		return appName + " Verification: " + args[0], []string{
			"Hello from " + appName + "!",
			" ",
			"Thanks for registering for our town. We hope you enjoy your time in this Liphium town and wish you an amazing stay.",
			"In case you haven't seen the verification code in the title of the email, here it is again: " + args[0],
			" ",
			"If you haven't registered for this town, please let the owner of this town know. Please do not contact Liphium support about it as we do not control this town. You can find out more at https://liphium.com.",
			" ",
			"Thank you for registering,",
			"The " + appName + " team.",
		}
	},
}
