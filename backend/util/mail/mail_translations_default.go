package mail

import (
	"node-backend/util"
	"os"
)

// English is the default for email translations yk
var enUsTranslations = map[string]func([]string) (string, []string){
	EmailVerification: func(args []string) (string, []string) {
		appName := os.Getenv(util.EnvAppName)
		return appName + " Email Verification Code", []string{
			"Hello from " + appName + "!",
			" ",
			"We are messaging you to verify your email and stuff, cause you registered an account! If this wasn't you, well then we might have a litte problem on our hands here. If this happened, please visit our website and report this issue to us.",
			"Well anyway, here's your verification code: " + args[0],
			" ",
			"Thank you for registering,",
			"The " + appName + " team.",
		}
	},
}
