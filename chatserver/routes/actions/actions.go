package remote_action_routes

import (
	"github.com/Liphium/station/chatserver/caching"
	conversation_actions "github.com/Liphium/station/chatserver/routes/actions/conversations"
	remote_event_channel "github.com/Liphium/station/chatserver/routes/actions/event_channel"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	message_actions "github.com/Liphium/station/chatserver/routes/actions/messages"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

// Setup the routes
func SetupRemoteActions(router fiber.Router) {

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
	router.Post("/ping", actionHandler(pingTest))
	router.Post("/negotiate", actionHandler(handleNegotiation))
	router.Post("/conv_subscribe", actionHandler(conversation_actions.HandleRemoteSubscription))
}

// Setup the event channel for nodes outside of the current one sending events for conversations on them
func SetupEventChannel(router fiber.Router) {
	router.Post("/send", remote_event_channel.HandleRemoteEvent)
}

// Setup all the actions that can be called from outside of the current node for a conversation on the current node
func SetupConversationActions(router fiber.Router) {

	// Actions for conversation management
	router.Post("/conv_activate", conversationHandler(conversation_actions.HandleTokenActivation))
	router.Post("/conv_promote", conversationHandler(conversation_actions.HandlePromoteToken))
	router.Post("/conv_demote", conversationHandler(conversation_actions.HandleDemoteToken))
	router.Post("/conv_read", conversationHandler(conversation_actions.HandleRead))
	router.Post("/conv_data", conversationHandler(conversation_actions.HandleGetData))
	router.Post("/conv_gen_token", conversationHandler(conversation_actions.HandleGenerateToken))
	router.Post("/conv_kick", conversationHandler(conversation_actions.HandleKick))
	router.Post("/conv_leave", conversationHandler(conversation_actions.HandleLeave))
	router.Post("/conv_st_res", conversationHandler(conversation_actions.HandleStatusResponse))

	// Actions for message management
	router.Post("/msg_delete", conversationHandler(message_actions.HandleDelete))
	router.Post("/msg_get", conversationHandler(message_actions.HandleGet))
	router.Post("/msg_list_after", conversationHandler(message_actions.HandleListAfter))
	router.Post("/msg_list_before", conversationHandler(message_actions.HandleListBefore))
	router.Post("/msg_send", conversationHandler(message_actions.HandleSend))
}

// Creates a new handler for the action based on its calling method
func actionHandler[T any](handler action_helpers.ActionHandlerFunc[T]) func(*fiber.Ctx) error {
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
		var req action_helpers.ConversationActionRequest[T]
		if err := integration.BodyParser(c, &req); err != nil {
			return integration.InvalidRequest(c, "action wasn't valid")
		}

		// Check if the token has been negotiated with before
		node, valid := tokenMap.Load(req.Token.ID)
		if !valid {
			return integration.InvalidRequest(c, "negotiation required")
		}

		// Add it to the locals
		c.Locals("node", node)

		// Check the conversation token
		token, err := caching.ValidateToken(req.Token.ID, req.Token.Token)
		if err != nil {

			// Delete the token from the negotiation map in case it is still in there
			tokenMap.Delete(token.ID)

			return integration.InvalidRequest(c, "token wasn't valid")
		}

		// Handle the action
		return handler(c, token, req.Data)
	}
}
