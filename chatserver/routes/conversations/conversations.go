package conversation_routes

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	message_routes "github.com/Liphium/station/chatserver/routes/conversations/message"
	"github.com/Liphium/station/chatserver/util"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
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

	// Begin a new transaction
	database.DBConn.Begin()
	defer func() {
		if err := recover(); err != nil {
			util.Log.Println("fatal error during transaction:", err)
			database.DBConn.Rollback()
		}
	}()

	// Get the conversation (in case it has changed)
	var currentVersion int64
	if err := database.DBConn.Model(&conversations.Conversation{}).Select("version").Where("id = ?", conversation.ID).Take(&currentVersion).Error; err != nil {
		database.DBConn.Rollback()
		return err
	}

	// Update the conversation
	conversation.Version = currentVersion + 1

	// Save the conversation
	if err := database.DBConn.Save(&conversation).Error; err != nil {
		database.DBConn.Rollback()
		return err
	}

	// Stop the transaction
	database.DBConn.Commit()
	return nil
}
