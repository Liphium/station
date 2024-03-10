package util

import (
	"log"
	"os"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

// General errors
const ErrorNode = "node.error"
const ErrorServer = "server.error"
const ErrorMail = "mail.error"

// Auth errors
const EmailInvalid = "email.invalid"
const CodeInvalid = "code.invalid"
const PasswordInvalid = "password.incorrect"
const EmailRegistered = "email.registered" // When it is already registered
const UsernameInvalid = "username.invalid"
const UsernameTaken = "username.taken"
const TagInvalid = "tag.invalid"
const InviteInvalid = "invite.invalid"

func DebugRouteError(c *fiber.Ctx, msg string) {
	if Testing {
		log.Println(c.Route().Path+":", msg)
	}
}

func SuccessfulRequest(c *fiber.Ctx) error {
	return ReturnJSON(c, fiber.Map{
		"success": true,
	})
}

func FailedRequest(c *fiber.Ctx, error string, err error) error {

	if LogErrors && err != nil {
		log.Println(c.Route().Path+":", err)
		if os.Getenv("SHOW_STACK") == "1" {
			debug.PrintStack()
		}
	}

	return ReturnJSON(c, fiber.Map{
		"success": false,
		"error":   error,
	})
}

func InvalidRequest(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusBadRequest)
}
