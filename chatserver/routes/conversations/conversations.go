package conversation_routes

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	message_routes "github.com/Liphium/station/chatserver/routes/conversations/message"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Unauthorized(router fiber.Router) {
	router.Post("/remote_activate", remoteActivate)
}

func Authorized(router fiber.Router) {
	router.Post("/open", openConversation)
	router.Post("/read", read)
	router.Post("/activate", activate)
	router.Post("/demote_token", demoteToken)
	router.Post("/promote_token", promoteToken)
	router.Post("/data", getData)
	router.Post("/generate_token", generateToken)
	router.Post("/kick_member", kickMember)
	router.Post("/leave", leaveConversation)
	router.Post("/change_data", changeData)

	router.Route("/message", message_routes.SetupRoutes)
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
