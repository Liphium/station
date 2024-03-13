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
	"github.com/Liphium/station/pipes/send"
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
	return len(r.Conversation) > 0 && len(r.Data) > 0 && len(r.Token) == util.ConversationTokenLength &&
		uint64(time.Now().UnixMilli())-r.Timestamp < 2000
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
		if member.TokenID == req.TokenID {
			found = true
		}
	}

	if !found {
		return integration.InvalidRequest(c, "member token wasn't found "+req.Token+" "+req.Conversation)
	}

	messageId := util.GenerateToken(32)
	certificate, err := conversations.GenerateCertificate(messageId, req.Conversation, req.TokenID)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	util.Log.Println(certificate)

	message := conversations.Message{
		ID:           messageId,
		Conversation: req.Conversation,
		Certificate:  certificate,
		Data:         req.Data,
		Sender:       req.TokenID,
		Creation:     int64(req.Timestamp),
		Edited:       false,
	}

	if err := database.DBConn.Create(&message).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("conversation = ? AND id = ?", req.Conversation, req.TokenID).Update("last_read", time.Now().UnixMilli()+1).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}
	token.LastRead = time.Now().UnixMilli() + 1
	caching.UpdateToken(token)

	adapters, nodes := caching.MembersToPipes(members)
	event := MessageEvent(message)

	send.Pipe(send.ProtocolWS, pipes.Message{
		Channel: pipes.Conversation(adapters, nodes),
		Event:   event,
	})

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
