package cluster

import "github.com/gofiber/fiber/v2"

func Setup(router fiber.Router) {
	router.Post("/add", addCluster)
	router.Post("/remove", removeCluster)
	router.Post("/list", listClusters)
}
