package remote_action_routes

import (
	"github.com/Liphium/station/chatserver/caching"
	conversation_actions "github.com/Liphium/station/chatserver/routes/actions/conversations"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

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
		if req["id"] == nil || req["token"] == nil || req["data"] == nil {
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

	// Conversation actions
	router.Post("/conv_activate", conversationHandler(conversation_actions.HandleTokenActivation))
	router.Post("/conv_promote", conversationHandler(conversation_actions.HandlePromoteToken))
	router.Post("/conv_demote", conversationHandler(conversation_actions.HandleDemoteToken))
	router.Post("/conv_read", conversationHandler(conversation_actions.HandleRead))
	router.Post("/conv_data", conversationHandler(conversation_actions.HandleGetData))
	router.Post("/conv_gen_token", conversationHandler(conversation_actions.HandleGenerateToken))
	router.Post("/conv_kick", conversationHandler(conversation_actions.HandleKick))
	router.Post("/conv_leave", conversationHandler(conversation_actions.HandleLeave))
}

// Creates a new handler for the action based on its calling method
func ActionHandler[T any](handler action_helpers.ActionHandlerFunc[T]) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Parse the action with the request generic
		var req action_helpers.RemoteActionRequest[T]
		if err := integration.BodyParser(c, &req); err != nil {
			return integration.InvalidRequest(c, "action wasn't valid")
		}

		// Add the remote action request data to the locals
		c.Locals("sender", req.Sender)

		// Handle the action
		return handler(c, req.Data)
	}
}

// Creates a new handler for the action based on its calling method
func conversationHandler[T any](handler action_helpers.ConversationActionHandlerFunc[T]) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Parse the action with the request generic
		var req action_helpers.RemoteActionRequest[action_helpers.ConversationActionRequest[T]]
		if err := integration.BodyParser(c, &req); err != nil {
			return integration.InvalidRequest(c, "action wasn't valid")
		}

		// Check the conversation token
		token, err := caching.ValidateToken(req.Data.Token.ID, req.Data.Token.Token)
		if err != nil {
			return integration.InvalidRequest(c, "token wasn't valid")
		}

		// Add the remote action request data to the locals
		c.Locals("sender", req.Sender)

		// Handle the action
		return handler(c, token, req.Data.Data)
	}
}
