package townhall_accounts

import "github.com/gofiber/fiber/v2"

func Authorized(router fiber.Router) {
	router.Post("/list", listAccounts)
	router.Post("/delete", deleteAccount)
	router.Post("/change_rank", changeRank)
}
