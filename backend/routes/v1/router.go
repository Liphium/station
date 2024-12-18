package routes_v1

import (
	"crypto/rsa"
	"encoding/base64"
	"os"

	"github.com/Liphium/station/backend/routes/v1/account"
	"github.com/Liphium/station/backend/routes/v1/node"
	townhall_routes "github.com/Liphium/station/backend/routes/v1/townhall"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/gofiber/fiber/v2"
)

func Router(router fiber.Router) {

	// Get default private and public key
	serverPublicKey, err := util.UnpackageRSAPublicKey(os.Getenv("TC_PUBLIC_KEY"))
	if err != nil {
		panic("Couldn't unpackage public key. Required for v1 API. Please set TC_PUBLIC_KEY in your environment variables or .env file. \n err: " + err.Error())
	}

	serverPrivateKey, err := util.UnpackageRSAPrivateKey(os.Getenv("TC_PRIVATE_KEY"))
	if err != nil {
		panic("Couldn't unpackage private key. Required for v1 API. Please set TC_PRIVATE_KEY in your environment variables or .env file. \n err: " + err.Error())
	}

	// Endpoint to get server public key (so no requirements apply yet)
	router.Post("/pub", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"pub":              util.PackageRSAPublicKey(serverPublicKey),
			"protocol_version": util.ProtocolVersion,
		})
	})

	// Use a middleware to make sure all the translations work properly
	router.Use(func(c *fiber.Ctx) error {

		// Set the locale for translations to work properly
		localeHeader, valid := c.GetReqHeaders()["Locale"]
		if valid {
			c.Locals("locale", localeHeader[0])
		}

		return c.Next()
	})

	// Unencrypted account routes
	router.Route("/v1/account", account.Unencrypted)
	router.Route("/v1", func(router fiber.Router) {
		encryptedRoutes(router, serverPublicKey, serverPrivateKey)
	})
}

func encryptedRoutes(router fiber.Router, serverPublicKey *rsa.PublicKey, serverPrivateKey *rsa.PrivateKey) {

	// Through Cloudflare Protection (Decryption method)
	router.Use(func(c *fiber.Ctx) error {

		// Check if the auth tag exists
		aesKeyEncoded, valid := c.GetReqHeaders()["Auth-Tag"]
		if !valid {
			util.Log.Println("no header")
			return c.SendStatus(fiber.StatusPreconditionFailed)
		}

		// Decode the auth tag
		aesKeyEncrypted, err := base64.StdEncoding.DecodeString(aesKeyEncoded[0])
		if err != nil {
			util.Log.Println("no decoding")
			return c.SendStatus(fiber.StatusPreconditionFailed)
		}

		// Get AES key from auth tag with server private key
		aesKey, err := util.DecryptRSA(serverPrivateKey, aesKeyEncrypted)
		if err != nil {
			return c.SendStatus(fiber.StatusPreconditionRequired)
		}

		// Decrypt the content of the request with the AES key
		decrypted, err := util.DecryptAES(aesKey, c.Body())
		if err != nil {
			return c.SendStatus(fiber.StatusNetworkAuthenticationRequired)
		}

		// Add some variables to work with the keys
		c.Locals(util.LocalsBody, decrypted)
		c.Locals(util.LocalsKey, aesKey)
		c.Locals(util.LocalsServerPub, serverPublicKey)
		c.Locals(util.LocalsServerPriv, serverPrivateKey)

		return c.Next()
	})

	// Unauthorized routes
	router.Route("/node", node.Unauthorized)
	router.Route("/account", account.Unauthorized)

	router.Route("/", authorizedRoutes)
}

func authorizedRoutes(router fiber.Router) {

	// Authorized by using a jwt token + verifiying it with the database
	router.Use(verify.AuthMiddleware())

	// Authorized routes
	router.Route("/account", account.Authorized)
	router.Route("/node", node.Authorized)
	router.Route("/townhall", townhall_routes.Authorized)
}
