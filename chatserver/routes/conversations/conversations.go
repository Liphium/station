package conversation_routes

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	conversation_actions "github.com/Liphium/station/chatserver/routes/actions/conversations"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	message_routes "github.com/Liphium/station/chatserver/routes/conversations/message"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Unauthorized(router fiber.Router) {
	router.Post("/remote_activate", remoteActivate)
}

func Authorized(router fiber.Router) {
	router.Post("/open", openConversation)
	router.Post("/read", read)
	router.Post("/activate", handler(conversation_actions.HandleTokenActivation))
	router.Post("/promote_token", handler(conversation_actions.HandlePromoteToken))
	router.Post("/promote_token", handler(conversation_actions.HandleDemoteToken))
	router.Post("/data", getData)
	router.Post("/generate_token", generateToken)
	router.Post("/kick_member", kickMember)
	router.Post("/leave", leaveConversation)
	router.Post("/change_data", changeData)

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

// This deletes all data related to a conversation
func deleteConversation(id string) error {
	if err := database.DBConn.Where("conversation = ?", id).Delete(&conversations.Message{}).Error; err != nil {
		return err
	}
	if err := database.DBConn.Where("conversation = ?", id).Delete(&conversations.ConversationToken{}).Error; err != nil {
		return err
	}
	if err := database.DBConn.Where("id = ?", id).Delete(&conversations.Conversation{}).Error; err != nil {
		return err
	}
	return nil
}

// This increments the version of the conversation by one in a transaction.
// Will also save the conversation.
func incrementConversationVersion(conversation conversations.Conversation) error {

	// Increment the version in a transaction
	err := database.DBConn.Transaction(func(tx *gorm.DB) error {

		// Get the current version (in case it has changed)
		var currentVersion int64
		if err := tx.Model(&conversations.Conversation{}).Select("version").Where("id = ?", conversation.ID).Take(&currentVersion).Error; err != nil {
			database.DBConn.Rollback()
			return err
		}

		// Update the conversation
		conversation.Version = currentVersion + 1

		// Save the conversation
		if err := tx.Save(&conversation).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
