package integration

import (
	"encoding/base64"
	"runtime/debug"
	"strings"

	"github.com/Liphium/station/main/localization"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func SuccessfulRequest(c *fiber.Ctx) error {
	return ReturnJSON(c, fiber.Map{
		"success": true,
	})
}

func FailedRequest(c *fiber.Ctx, message localization.Translations, err error) error {

	// Print error if it isn't nil
	if err != nil {
		Log.Println(c.Route().Name + " ERROR: " + message[localization.DefaultLocale] + ":" + err.Error())
		debug.PrintStack()
	}

	return ReturnJSON(c, fiber.Map{
		"success": false,
		"error":   Translate(c, message),
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

// Translate any message on a request
func Translate(c *fiber.Ctx, message localization.Translations) string {
	locale := c.Locals("locale")
	if locale == nil {
		locale = localization.DefaultLocale
	}
	msg, valid := message[strings.ToLower(locale.(string))]
	if !valid {
		msg = message[localization.DefaultLocale]
	}
	return msg
}

// Return encrypted json
func ReturnJSON(c *fiber.Ctx, data interface{}) error {

	encoded, err := sonic.Marshal(data)
	if err != nil {
		return FailedRequest(c, localization.ErrorServer, err)
	}

	if c.Locals("key") == nil {
		return c.Send(encoded)
	}
	encrypted, err := EncryptAES(c.Locals("key").([]byte), encoded)
	if err != nil {
		return FailedRequest(c, localization.ErrorServer, err)
	}

	return c.Send(encrypted)
}

func ThroughCloudflareMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		// Get the AES encryption key from the Auth-Tag header
		aesKeyEncoded, valid := c.GetReqHeaders()["Auth-Tag"]
		if !valid {
			Log.Println("no header")
			return c.SendStatus(fiber.StatusPreconditionFailed)
		}
		aesKeyEncrypted, err := base64.StdEncoding.DecodeString(aesKeyEncoded[0])
		if err != nil {
			Log.Println("no decoding")
			return c.SendStatus(fiber.StatusPreconditionFailed)
		}

		// Decrypt the AES key using the private key of this node
		aesKey, err := DecryptRSA(NodePrivateKey, aesKeyEncrypted)
		if err != nil {
			return c.SendStatus(fiber.StatusPreconditionRequired)
		}

		// Decrypt the request body using the key attached to the Auth-Tag header
		decrypted, err := DecryptAES(aesKey, c.Body())
		if err != nil {
			return c.SendStatus(fiber.StatusNetworkAuthenticationRequired)
		}

		// Set some variables for use when sending back the response
		c.Locals("body", decrypted)
		c.Locals("key", aesKey)

		// Go to the next middleware/handler
		return c.Next()
	}
}
