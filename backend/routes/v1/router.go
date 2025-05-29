package routes_v1

import (
	"github.com/Liphium/station/backend/routes/v1/account"
	"github.com/Liphium/station/backend/routes/v1/node"
	townhall_routes "github.com/Liphium/station/backend/routes/v1/townhall"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/gofiber/fiber/v2"
)

func Router(router fiber.Router) {

	// Endpoint to get server public key (so no requirements apply yet)
	router.Post("/about", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
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
	router.Route("/v1", encryptedRoutes)
}

func encryptedRoutes(router fiber.Router) {

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
