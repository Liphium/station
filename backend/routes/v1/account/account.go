package account

import (
	auth_routes "github.com/Liphium/station/backend/routes/v1/account/auth"
	friends2_routes "github.com/Liphium/station/backend/routes/v1/account/friends2"
	invite_routes "github.com/Liphium/station/backend/routes/v1/account/invite"
	"github.com/Liphium/station/backend/routes/v1/account/keys"
	"github.com/Liphium/station/backend/routes/v1/account/profile"
	"github.com/Liphium/station/backend/routes/v1/account/rank"
	settings_routes "github.com/Liphium/station/backend/routes/v1/account/settings"
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Route("/auth", auth_routes.Unauthorized)
	router.Route("/keys", keys.Unauthorized)
	router.Route("/friends", friends2_routes.Unauthorized)
	// router.Route("/file_info", files.Unauthorized) // This is this way because of routing issues with fiber (can't blame them though, this is way too complex)
	router.Route("/rank", rank.Unauthorized)
	// router.Route("/stored_actions", stored_actions.Unauthorized)

	router.Route("/profile", profile.Unauthorized)
	router.Post("/get", getAccount)
	// router.Post("/get_node", getAccountNode)
}

func Authorized(router fiber.Router) {
	router.Route("/keys", keys.Authorized)
	// router.Route("/stored_actions", stored_actions.Authorized)
	router.Route("/friends", friends2_routes.Authorized)
	// router.Route("/vault", vault.Authorized)
	router.Route("/profile", profile.Authorized)
	router.Route("/invite", invite_routes.Authorized)
	// router.Route("/files", files.Authorized)
	router.Route("/settings", settings_routes.Authorized)

	router.Post("/me", me)
	router.Post("/get_name", getAccountByUsername)
}
