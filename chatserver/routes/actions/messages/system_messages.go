package message_actions

import (
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/bytedance/sonic"
)

const systemSender = "6969@liphium.com"

// Stored system messages
const GroupNewAdmin = "group.new_admin"
const GroupRankChange = "group.rank_change"
const GroupMemberJoin = "group.member_join"
const GroupMemberKick = "group.member_kick"
const GroupMemberInvite = "group.member_invite"
const GroupMemberLeave = "group.member_leave"
const ConversationEdited = "conv.edited"

// Not stored system messages
const ConversationVersionUpdate = "conv.update"
const DeletedMessage = "msg.deleted"
const ConversationKick = "conv.kicked"

// Send a system message that is stored in the database
func SendSystemMessage(conversation string, content string, attachments []string) error {

	contentJson, err := sonic.MarshalString(map[string]interface{}{
		"c": content,
		"a": attachments,
	})
	if err != nil {
		return err
	}

	message := conversations.Message{
		Conversation: conversation,
		Data:         contentJson,
		Sender:       systemSender,
		Creation:     time.Now().UnixMilli(),
		Edited:       false,
	}

	// Save message to the database
	if err := database.DBConn.Create(&message).Error; err != nil {
		return err
	}

	// Load members
	members, err := caching.LoadMembers(conversation)
	if err != nil {
		return err
	}

	// Send the event to all the members
	event := MessageEvent(message)
	return caching.SendEventToMembers(members, event)
}

// Send a system message that isn't stored in the database
func SendNotStoredSystemMessage(conversation string, content string, attachments []string) error {

	// Generate the content for the message
	contentJson, err := sonic.MarshalString(map[string]interface{}{
		"c": content,
		"a": attachments,
	})
	if err != nil {
		return err
	}

	// Create the message
	message := conversations.Message{
		Conversation: conversation,
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

	// Send the event to the members
	event := MessageEvent(message)
	return caching.SendEventToMembers(members, event)
}

func AttachAccount(encrypted string) string {
	return "a:" + encrypted
}
