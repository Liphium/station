package townhall_routes

import (
	townhall_accounts "github.com/Liphium/station/backend/routes/v1/townhall/accounts"
	townhall_settings "github.com/Liphium/station/backend/routes/v1/townhall/settings"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {

	// Add a middleware to prevent non-admin accounts from accessing these endpoints
	router.Use(func(c *fiber.Ctx) error {

		// Deny access if the user doesn't have admin permissions
		if !verify.InfoLocals(c).HasPermission(verify.PermissionAdmin) {
			return util.InvalidRequest(c)
		}

		return c.Next()
	})

	router.Route("/accounts", townhall_accounts.Authorized)
	router.Route("/settings", townhall_settings.Authorized)
}
