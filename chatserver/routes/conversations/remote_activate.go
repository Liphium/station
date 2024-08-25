package conversation_routes

import (
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type remoteActivateRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Node  string `json:"node"` // The node the request came from
}

// Route: /conversations/remote_activate
func remoteActivate(c *fiber.Ctx) error {

	// Parse the request
	var req remoteActivateRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "request was invalid")
	}

	//

	return integration.SuccessfulRequest(c)
}
