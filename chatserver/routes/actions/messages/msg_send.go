package message_actions

import (
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/gofiber/fiber/v2"
)

type MessageSendAction struct {
	Timestamp uint64 `json:"timestamp"`
	Data      string `json:"data"`
}

// Action: msg_send
func HandleSend(c *fiber.Ctx, token conversations.ConversationToken, action MessageSendAction) error {

	// Validate request
	if len(action.Data) == 0 {
		return integration.InvalidRequest(c, "request is invalid")
	}

	// Check if message is too big
	if conversations.CheckSize(action.Data) {
		return integration.FailedRequest(c, localization.ErrorMessageTooLong, nil)
	}

	// Generate an id and certificate for the message
	messageId := util.GenerateToken(32)
	certificate, err := conversations.GenerateCertificate(messageId, token.Conversation, token.ID)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	message := conversations.Message{
		ID:           messageId,
		Conversation: token.Conversation,
		Certificate:  certificate,
		Data:         action.Data,
		Sender:       token.ID,
		Creation:     int64(action.Timestamp),
		Edited:       false,
	}

	// Save the message to the database
	if err := database.DBConn.Create(&message).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Update the read state to prevent the message sender from being notified about the message
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("conversation = ? AND id = ?", token.Conversation, token.ID).Update("last_read", time.Now().UnixMilli()+1).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Load the members of the conversation
	members, err := caching.LoadMembers(token.Conversation)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send the message to everyone
	event := MessageEvent(message)
	if err := caching.SendEventToMembers(members, event); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"message": message,
	})
}

func MessageEvent(message conversations.Message) pipes.Event {
	return pipes.Event{
		Name: "conv_msg",
		Data: map[string]interface{}{
			"conv": message.Conversation,
			"msg":  message,
		},
	}
}
