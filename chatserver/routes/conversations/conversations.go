package conversation_routes

import (
	conversation_actions "github.com/Liphium/station/chatserver/routes/actions/conversations"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	message_routes "github.com/Liphium/station/chatserver/routes/conversations/message"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {
	router.Post("/open", openConversation)
	router.Post("/read", handler(conversation_actions.HandleRead))
	router.Post("/activate", handler(conversation_actions.HandleTokenActivation))
	router.Post("/promote_token", handler(conversation_actions.HandlePromoteToken))
	router.Post("/promote_token", handler(conversation_actions.HandleDemoteToken))
	router.Post("/data", handler(conversation_actions.HandleGetData))
	router.Post("/generate_token", handler(conversation_actions.HandleGenerateToken))
	router.Post("/kick_member", handler(conversation_actions.HandleKick))
	router.Post("/leave", handler(conversation_actions.HandleLeave))
	router.Post("/change_data", handler(conversation_actions.HandleChangeData))

	router.Route("/message", message_routes.SetupRoutes)
}

// Create a normal endpoint from an action handler
func handler[T any](handler action_helpers.ActionHandlerFunc[T]) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		// Parse the request
		var req T
		if err := integration.BodyParser(c, &req); err != nil {
			return integration.InvalidRequest(c, "request was invalid")
		}

		// Let the action handle the request
		return handler(c, req)
	}
}
