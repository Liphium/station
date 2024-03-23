package integration

import (
	"runtime/debug"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

const ErrorServer = "server.error"

func SuccessfulRequest(c *fiber.Ctx) error {
	return ReturnJSON(c, fiber.Map{
		"success": true,
	})
}

func FailedRequest(c *fiber.Ctx, message string, err error) error {

	// Print error if it isn't nil
	if err != nil {
		Log.Println(c.Route().Name + " ERROR: " + message + ":" + err.Error())
		debug.PrintStack()
	}

	return ReturnJSON(c, fiber.Map{
		"success": false,
		"error":   message,
	})
}

func InvalidRequest(c *fiber.Ctx, message string) error {
	Log.Println(c.Route().Name + " request is invalid. msg: " + message)
	debug.PrintStack()
	return c.SendStatus(fiber.StatusBadRequest)
}

// Parse encrypted json
func BodyParser(c *fiber.Ctx, data interface{}) error {
	return sonic.Unmarshal(c.Locals("body").([]byte), data)
}

// Return encrypted json
func ReturnJSON(c *fiber.Ctx, data interface{}) error {

	encoded, err := sonic.Marshal(data)
	if err != nil {
		return FailedRequest(c, ErrorServer, err)
	}

	if c.Locals("key") == nil {
		return c.Send(encoded)
	}
	encrypted, err := EncryptAES(c.Locals("key").([]byte), encoded)
	if err != nil {
		return FailedRequest(c, ErrorServer, err)
	}

	return c.Send(encrypted)
}
