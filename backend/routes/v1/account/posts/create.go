package post_routes

import (
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/posts/create
func createPost(c *fiber.Ctx) error {
	return util.SuccessfulRequest(c)
}
