package conversation_routes

import (
	conversation_actions "github.com/Liphium/station/chatserver/routes/actions/conversations"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	message_actions "github.com/Liphium/station/chatserver/routes/actions/messages"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {

	// Setup all conversation routes
	router.Post("/open", openConversation)
	router.Post("/read", action_helpers.CreateConversationEndpoint(conversation_actions.HandleRead, "conv_read"))
	router.Post("/timestamp", action_helpers.CreateConversationEndpoint(conversation_actions.HandleTimestamp, "conv_timestamp"))
	router.Post("/activate", action_helpers.CreateConversationEndpoint(conversation_actions.HandleTokenActivation, "conv_activate"))
	router.Post("/promote_token", action_helpers.CreateConversationEndpoint(conversation_actions.HandlePromoteToken, "conv_promote"))
	router.Post("/demote_token", action_helpers.CreateConversationEndpoint(conversation_actions.HandleDemoteToken, "conv_demote"))
	router.Post("/data", action_helpers.CreateConversationEndpoint(conversation_actions.HandleGetData, "conv_data"))
	router.Post("/set_data", action_helpers.CreateConversationEndpoint(conversation_actions.HandleSetData, "conv_set_data"))
	router.Post("/generate_token", action_helpers.CreateConversationEndpoint(conversation_actions.HandleGenerateToken, "conv_gen_token"))
	router.Post("/kick_member", action_helpers.CreateConversationEndpoint(conversation_actions.HandleKick, "conv_kick"))
	router.Post("/leave", action_helpers.CreateConversationEndpoint(conversation_actions.HandleLeave, "conv_leave"))
	router.Post("/answer_status", action_helpers.CreateConversationEndpoint(conversation_actions.HandleStatusResponse, "conv_st_res"))

	// Setup all message routes
	router.Post("/message/send", action_helpers.CreateConversationEndpoint(message_actions.HandleSend, "msg_send"))
	router.Post("/message/delete", action_helpers.CreateConversationEndpoint(message_actions.HandleDelete, "msg_delete"))
	router.Post("/message/list_after", action_helpers.CreateConversationEndpoint(message_actions.HandleListAfter, "msg_list_after"))
	router.Post("/message/list_before", action_helpers.CreateConversationEndpoint(message_actions.HandleListBefore, "msg_list_before"))
	router.Post("/message/get", action_helpers.CreateConversationEndpoint(message_actions.HandleGet, "msg_get"))
}
