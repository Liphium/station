package message_routes

import (
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
	"github.com/gofiber/fiber/v2"
)

type MessageSendRequest struct {
	Conversation string `json:"conversation"`
	TokenID      string `json:"token_id"`
	Token        string `json:"token"`
	Timestamp    uint64 `json:"timestamp"`
	Data         string `json:"data"`
}

func (r *MessageSendRequest) Validate() bool {
	return len(r.Conversation) > 0 && len(r.Data) > 0 && len(r.Token) == util.ConversationTokenLength
}

// Route: /conversations/message/send
func sendMessage(c *fiber.Ctx) error {

	var req MessageSendRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, err.Error())
	}

	// Validate request
	if !req.Validate() {
		return integration.InvalidRequest(c, "request is invalid")
	}

	if conversations.CheckSize(req.Data) {
		return integration.FailedRequest(c, "too.big", nil)
	}

	// Validate conversation token
	token, err := caching.ValidateToken(req.TokenID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, "token id is invalid")
	}

	// Load members
	members, err := caching.LoadMembers(req.Conversation)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, nil)
	}

	found := false
	for _, member := range members {
		if member.TokenID == token.ID {
			found = true
		}
	}

	if !found {
		return integration.InvalidRequest(c, "member token wasn't found "+req.Token+" "+req.Conversation)
	}

	// Generate an id and certificate for the message
	messageId := util.GenerateToken(32)
	certificate, err := conversations.GenerateCertificate(messageId, req.Conversation, token.ID)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	message := conversations.Message{
		ID:           messageId,
		Conversation: req.Conversation,
		Certificate:  certificate,
		Data:         req.Data,
		Sender:       token.ID,
		Creation:     int64(req.Timestamp),
		Edited:       false,
	}

	// Save the message to the database
	if err := database.DBConn.Create(&message).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Update the read state to prevent the message sender from being notified about the message
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("conversation = ? AND id = ?", req.Conversation, req.TokenID).Update("last_read", time.Now().UnixMilli()+1).Error; err != nil {
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
