package message_routes

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Request for deleting a message
type deleteMessageRequest struct {
	TokenID     string `json:"id"`          // Conversation token id
	Token       string `json:"token"`       // Conversation token (token)
	Certificate string `json:"certificate"` // Message certificate
}

// Route: /conversations/message/delete
func deleteMessage(c *fiber.Ctx) error {

	// Parse request
	var req deleteMessageRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get conversation token
	token, err := caching.ValidateToken(req.TokenID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, "invalid conversation token")
	}

	// Get claims from message certificate
	claims, valid := conversations.GetCertificateClaims(req.Certificate)
	if !valid {
		return integration.InvalidRequest(c, "invalid certificate claims")
	}

	util.Log.Println(claims)

	// Check if certificate is valid for the provided conversation token
	util.Log.Println("message:", claims.Message, claims.Message)
	util.Log.Println("conv:", claims.Conversation, token.Conversation)
	util.Log.Println("sender:", claims.Sender, token.ID)
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
