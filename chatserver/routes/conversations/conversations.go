package conversation_routes

import (
	"strings"

	"github.com/Liphium/station/chatserver/caching"
	conversation_actions "github.com/Liphium/station/chatserver/routes/actions/conversations"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	message_routes "github.com/Liphium/station/chatserver/routes/conversations/message"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {
	router.Post("/open", openConversation)
	router.Post("/read", handler(conversation_actions.HandleRead, "conv_read"))
	router.Post("/activate", handler(conversation_actions.HandleTokenActivation, "conv_activate"))
	router.Post("/promote_token", handler(conversation_actions.HandlePromoteToken, "conv_promote"))
	router.Post("/demote_token", handler(conversation_actions.HandleDemoteToken, "conv_demote"))
	router.Post("/data", handler(conversation_actions.HandleGetData, "conv_data"))
	router.Post("/generate_token", handler(conversation_actions.HandleGenerateToken, "conv_gen_token"))
	router.Post("/kick_member", handler(conversation_actions.HandleKick, "conv_kick"))
	router.Post("/leave", handler(conversation_actions.HandleLeave, "conv_leave"))
	router.Post("/change_data", handler(conversation_actions.HandleChangeData, "conv_data"))

	router.Route("/message", message_routes.SetupRoutes)
}

// Create a normal endpoint from an conversation action handler
func handler[T any](handler action_helpers.ConversationActionHandlerFunc[T], action string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		// Parse the request
		var req action_helpers.ConversationActionRequest[T]
		if err := integration.BodyParser(c, &req); err != nil {
			return integration.InvalidRequest(c, "request was invalid")
		}

		// Parse the conversation to extract the address
		args := strings.Split(req.Token.Conversation, "@")
		if len(args) != 2 {
			return integration.InvalidRequest(c, "conversation id is invalid")
		}

		// If the address isn't the current instance, send a remote action
		if args[1] != integration.BasePath {

			// Send a remote action to the other instance
			res, err := integration.PostRequestTC(args[1], "actions/"+action, fiber.Map{
				"app_tag": integration.AppTagChatNode,
				"sender":  caching.CSNode.SL,
				"action":  action,
				"data":    req,
			})
			if err != nil {
				return integration.FailedRequest(c, localization.ErrorServer, err)
			}

			// Return the response to the client
			return integration.ReturnJSON(c, res)
		}

		// Validate the token
		token, err := caching.ValidateToken(req.Token.ID, req.Token.Token)
		if err != nil {
			return integration.InvalidRequest(c, "conversation token was valid")
		}

		// Let the action handle the request
		return handler(c, token, req.Data)
	}
}
