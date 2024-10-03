package townhall_routes

import (
	townhall_accounts "github.com/Liphium/station/backend/routes/v1/townhall/accounts"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {

	// Add a middleware to prevent non-admin accounts from accessing these endpoints
	router.Use(func(c *fiber.Ctx) error {

		// Deny access if the user doesn't have admin permissions
		if !util.Permission(c, util.PermissionAdmin) {
			return util.InvalidRequest(c)
		}

		return c.Next()
	})

	router.Route("/accounts", townhall_accounts.Authorized)
}
