package account

import (
	"node-backend/routes/v1/account/files"
	"node-backend/routes/v1/account/friends"
	invite_routes "node-backend/routes/v1/account/invite"
	"node-backend/routes/v1/account/keys"
	"node-backend/routes/v1/account/profile"
	"node-backend/routes/v1/account/rank"
	settings_routes "node-backend/routes/v1/account/settings"
	"node-backend/routes/v1/account/stored_actions"
	"node-backend/routes/v1/account/vault"

	"github.com/gofiber/fiber/v2"
)

func Unencrypted(router fiber.Router) {
	router.Route("/files", files.Unencrypted)
}

func Unauthorized(router fiber.Router) {
	router.Route("/rank", rank.Unauthorized)
	router.Route("/stored_actions", stored_actions.Unauthorized)
}

func Authorized(router fiber.Router) {
	router.Route("/keys", keys.Authorized)
	router.Route("/stored_actions", stored_actions.Authorized)
	router.Route("/friends", friends.Authorized)
	router.Route("/vault", vault.Authorized)
	router.Route("/profile", profile.Authorized)
	router.Route("/invite", invite_routes.Authorized)
	router.Route("/files", files.Authorized)
	router.Route("/settings", settings_routes.Authorized)

	router.Post("/me", me)
	router.Post("/get", getAccount)
}
