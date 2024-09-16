package message_actions

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: msg_delete
func HandleDelete(c *fiber.Ctx, token conversations.ConversationToken, messageId string) error {

	// Get the message
	var message conversations.Message
	if err := database.DBConn.Where("id = ?", messageId).Take(&message).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorMessageAlreadyDeleted, err)
	}

	// Make sure the deleter is the sender
	if message.Sender != token.ID {
		return integration.FailedRequest(c, localization.ErrorMessageDeleteNoPermission, nil)
	}

	// Delete the message in the database
	if err := database.DBConn.Where("id = ?", messageId).Delete(&conversations.Message{}).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send a system message to delete the message on all clients who are storing it
	if err := SendNotStoredSystemMessage(message.Conversation, DeletedMessage, []string{messageId}); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
