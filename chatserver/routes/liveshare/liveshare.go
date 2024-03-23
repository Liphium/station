package liveshare_routes

import "github.com/gofiber/fiber/v2"

func Unencrypted(router fiber.Router) {
	router.Post("/received", receivedChunk)
	router.Get("/download", downloadChunk)
	router.Post("/subscribe", subscribeToLiveshare)
	router.Post("/info", transactionInfo)
}

func Authorized(router fiber.Router) {
	router.Post("/upload", sendFilePart)
}
