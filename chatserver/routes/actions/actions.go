package remote_action_routes

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

// Generic struct to wrap the json with any additional data for an action
type remoteActionRequest[T any] struct {
	ID     string `json:"id"`
	Token  string `json:"token"`
	Action string `json:"action"`
	Data   T      `json:"data"`
}

// Setup the routes
func Unauthorized(router fiber.Router) {

	// Inject a middleware that checks the node token and id in the body
	router.Use(func(c *fiber.Ctx) error {

		// Parse the request
		var req map[string]interface{}
		if err := integration.BodyParser(c, &req); err != nil {
			return integration.InvalidRequest(c, "request is invalid")
		}

		// Check if the required data is existent
		if req["id"] == nil || req["token"] == nil || req["action"] == nil || req["data"] == nil {
			return integration.InvalidRequest(c, "request doesn't contain everything")
		}

		// Check if the data is valid
		if req["id"] != caching.CSNode.ID || req["token"] != caching.CSNode.Token {
			return integration.FailedRequest(c, localization.InvalidCredentials, nil)
		}

		return c.Next()
	})

	// All the actions
	router.Post("/ping", pingTest)
}
