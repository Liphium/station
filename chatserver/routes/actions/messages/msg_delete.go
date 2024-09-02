package message_actions

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Action: msg_delete
func HandleDelete(c *fiber.Ctx, token conversations.ConversationToken, certificate string) error {

	// Get claims from message certificate
	claims, valid := conversations.GetCertificateClaims(certificate)
	if !valid {
		return integration.InvalidRequest(c, "invalid certificate claims")
	}

	// Check if certificate is valid for the provided conversation token
	if !claims.Valid(claims.Message, token.Conversation, token.ID) {
		return integration.InvalidRequest(c, "no permssion to delete message")
	}

	// Check if there was already a deletion request
	contentJson, err := sonic.MarshalString(map[string]interface{}{
		"c": DeletedMessage,
		"a": []string{claims.Message},
	})
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}
	var justHereForNoNilPointer conversations.Message
	if err := database.DBConn.Where("data = ? AND conversation = ?", contentJson, claims.Conversation).Select("id").Take(&justHereForNoNilPointer).Error; err == nil {
		return integration.FailedRequest(c, "already.deleted", nil)
	}

	// Delete the message in the database
	if err := database.DBConn.Where("id = ?", claims.Message).Delete(&conversations.Message{}).Error; err != nil && err != gorm.ErrRecordNotFound {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send a system message to delete the message on all clients who are storing it
	if err := SendNotStoredSystemMessage(claims.Conversation, DeletedMessage, []string{claims.Message}); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
