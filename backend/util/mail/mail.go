package mail

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

// Email identifiers
const EmailVerification = "verification"

// Language -> Email Identifier -> Builder function
var emailBuilders = map[string]map[string]func([]string) (string, []string){
	"en_us": enUsTranslations,
	"de_de": deDeTranslations,
}

const defaultLocale = "en_us"

// Send an email
//
// email = address you wanna send to, locale = language of the email, name = email identifier (e.g. "verification"), args = arguments to the builder function
func SendEmail(email string, locale string, name string, args ...string) error {

	// Generate message
	translation, valid := emailBuilders[locale]
	if !valid {
		translation = emailBuilders[defaultLocale]
	}
	subject, body := translation[name](args)
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r%s", subject, strings.Join(body, "\n")))

	// Authenticate using the provided credentials
	auth := smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PW"), os.Getenv("SMTP_SERVER"))

	// Send the email
	err := smtp.SendMail(
		os.Getenv("SMTP_SERVER")+":"+os.Getenv("SMTP_PORT"),
		auth,
		os.Getenv("SMTP_FROM"),
		[]string{email},
		msg,
	)
	return err
}
