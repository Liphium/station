package townhall_settings

import "github.com/gofiber/fiber/v2"

func Authorized(router fiber.Router) {
	router.Post("/categories", getCategories)
	router.Post("/set_int", setIntegerSetting)
	router.Post("/set_bool", setBooleanSetting)

	// All settings return endpoints
	router.Post("/files", fileSettings)
	router.Post("/decentralization", decentralizationSettings)
}
