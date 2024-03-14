package message_routes

import (
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/pipes"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

const systemSender = "6969"

func SetupRoutes(router fiber.Router) {
	router.Post("/send", sendMessage)
	router.Post("/delete", deleteMessage)
}

// System messages
const DeletedMessage = "msg.deleted"
const GroupNewAdmin = "group.new_admin"
const GroupRankChange = "group.rank_change"
const GroupMemberJoin = "group.member_join"
const GroupMemberKick = "group.member_kick"
const GroupMemberInvite = "group.member_invite"
const GroupMemberLeave = "group.member_leave"

// Message not stored, but sent to just disconnect one person
const ConversationKick = "conv.kicked"

func SendSystemMessage(conversation string, content string, attachments []string) error {

	contentJson, err := sonic.MarshalString(map[string]interface{}{
		"c": content,
		"a": attachments,
	})
	if err != nil {
		return err
	}

	messageId := util.GenerateToken(32)
	message := conversations.Message{
		ID:           messageId,
		Conversation: conversation,
		Certificate:  "",
		Data:         contentJson,
		Sender:       systemSender,
		Creation:     time.Now().UnixMilli(),
		Edited:       false,
	}

	if err := database.DBConn.Create(&message).Error; err != nil {
		return err
	}

	// Load members
	members, err := caching.LoadMembers(conversation)
	if err != nil {
		return err
	}
	adapters, nodes := caching.MembersToPipes(members)

	event := MessageEvent(message)
	err = caching.Node.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.Conversation(adapters, nodes),
		Event:   event,
	})
	if err != nil {
		return err
	}

	return nil
}

func SendNotStoredSystemMessage(conversation string, content string, attachments []string) error {

	contentJson, err := sonic.MarshalString(map[string]interface{}{
		"c": content,
		"a": attachments,
	})
	if err != nil {
		return err
	}

	messageId := util.GenerateToken(32)
	message := conversations.Message{
		ID:           messageId,
		Conversation: conversation,
		Certificate:  "",
		Data:         contentJson,
		Sender:       systemSender,
		Creation:     time.Now().UnixMilli(),
		Edited:       false,
	}

	// Load members
	members, err := caching.LoadMembers(conversation)
	if err != nil {
		return err
	}
	adapters, nodes := caching.MembersToPipes(members)

	event := MessageEvent(message)
	err = caching.Node.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.Conversation(adapters, nodes),
		Event:   event,
	})
	if err != nil {
		return err
	}

	return nil
}

func AttachAccount(encrypted string) string {
	return "a:" + encrypted
}
