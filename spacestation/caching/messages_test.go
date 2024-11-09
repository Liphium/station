package caching

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddMessage(t *testing.T) {
	roomID := "room1"
	message := Message{
		ID:           "msg1",
		Conversation: roomID,
		Creation:     time.Now().Unix(),
		Data:         "Hello, world!",
		Edited:       false,
		Sender:       "user1",
	}

	messageMap.Store(roomID, &MessageSink{Mutex: &sync.Mutex{}, Messages: []Message{}})
	err := AddMessage(roomID, message)
	assert.NoError(t, err, "Adding a message should not return an error")

	// Verify message added
	sink, _ := messageMap.Load(roomID)
	msgSink := sink.(*MessageSink)
	assert.Equal(t, 1, len(msgSink.Messages), "There should be one message in the sink")
	assert.Equal(t, "msg1", msgSink.Messages[0].ID, "The message ID should match")
}

func TestGetMessagesBefore(t *testing.T) {
	roomID := "room1"
	messageMap.Store(roomID, &MessageSink{
		Mutex: &sync.Mutex{},
		Messages: []Message{
			{ID: "1", Creation: 10000},
			{ID: "2", Creation: 20000},
			{ID: "3", Creation: 30000},
			{ID: "4", Creation: 40000},
		},
	})

	messages, err := GetMessagesBefore(roomID, 30000)
	assert.NoError(t, err, "Getting messages before should not return an error")
	assert.Equal(t, 2, len(messages), "Should return two messages before the specified time")
	assert.Equal(t, "1", messages[0].ID, "First message ID should match")
	assert.Equal(t, "2", messages[1].ID, "Second message ID should match")
}

func TestGetMessagesBeforeBinarySearch(t *testing.T) {
	roomID := "room1"
	messageMap.Store(roomID, &MessageSink{
		Mutex: &sync.Mutex{},
		Messages: []Message{
			{ID: "1", Creation: 10000},
			{ID: "2", Creation: 20000},
			{ID: "3", Creation: 30000},
			{ID: "4", Creation: 40000},
			{ID: "5", Creation: 50000},
			{ID: "6", Creation: 60000},
			{ID: "7", Creation: 70000},
			{ID: "8", Creation: 80000},
			{ID: "9", Creation: 90000},
			{ID: "10", Creation: 100000},
			{ID: "11", Creation: 110000},
			{ID: "12", Creation: 120000},
			{ID: "13", Creation: 130000},
			{ID: "14", Creation: 140000},
			{ID: "15", Creation: 150000},
			{ID: "16", Creation: 160000},
		},
	})

	// Make sure the thing only returns 10
	messages, err := GetMessagesBefore(roomID, 140000)
	assert.NoError(t, err, "Getting messages before should not return an error")
	assert.Equal(t, 10, len(messages), "Should return two messages the 10 messages before the specified time.")
	assert.Equal(t, "13", messages[0].ID, "First message ID should match, actual messages:", messages)
	assert.Equal(t, "12", messages[1].ID, "Second message ID should match, actual messages:", messages)

	messages, err = GetMessagesBefore(roomID, 200000)
	assert.NoError(t, err, "Getting messages before should not return an error")
	assert.Equal(t, 10, len(messages), "Should return two messages the 10 messages before the specified time.")
	assert.Equal(t, "16", messages[0].ID, "First message ID should match, actual messages:", messages)
	assert.Equal(t, "15", messages[1].ID, "Second message ID should match, actual messages:", messages)

	// Make sure the thing doesn't return empty messages beyond the limit of 10
	messages, err = GetMessagesBefore(roomID, 30000)
	assert.NoError(t, err, "Getting messages before should not return an error")
	assert.Equal(t, 2, len(messages), "Should return two messages before the specified time.")
	assert.Equal(t, "2", messages[0].ID, "First message ID should match, actual messages:", messages)
	assert.Equal(t, "1", messages[1].ID, "Second message ID should match, actual messages:", messages)
}

func TestGetMessagesAfter(t *testing.T) {
	roomID := "room2"
	messageMap.Store(roomID, &MessageSink{
		Mutex: &sync.Mutex{},
		Messages: []Message{
			{ID: "1", Creation: 10000},
			{ID: "2", Creation: 20000},
			{ID: "3", Creation: 30000},
			{ID: "4", Creation: 40000},
		},
	})

	messages, err := GetMessagesAfter(roomID, 20000)
	assert.NoError(t, err, "Getting messages after should not return an error")
	assert.Equal(t, 2, len(messages), "Should return two messages after the specified time")
	assert.Equal(t, "3", messages[0].ID, "Third message ID should match")
	assert.Equal(t, "4", messages[1].ID, "Fourth message ID should match")
}

func TestGetMessagesAfterBinarySearch(t *testing.T) {
	roomID := "room2"
	messageMap.Store(roomID, &MessageSink{
		Mutex: &sync.Mutex{},
		Messages: []Message{
			{ID: "1", Creation: 10000},
			{ID: "2", Creation: 20000},
			{ID: "3", Creation: 30000},
			{ID: "4", Creation: 40000},
			{ID: "5", Creation: 50000},
			{ID: "6", Creation: 60000},
			{ID: "7", Creation: 70000},
			{ID: "8", Creation: 80000},
			{ID: "9", Creation: 90000},
			{ID: "10", Creation: 100000},
			{ID: "11", Creation: 110000},
			{ID: "12", Creation: 120000},
			{ID: "13", Creation: 130000},
			{ID: "14", Creation: 140000},
			{ID: "15", Creation: 150000},
			{ID: "16", Creation: 160000},
		},
	})

	// Make sure the thing only returns 10
	messages, err := GetMessagesAfter(roomID, 30000)
	assert.NoError(t, err, "Getting messages after should not return an error.")
	assert.Equal(t, 10, len(messages), "Should return two messages the 10 messages after the specified time.")
	assert.Equal(t, "4", messages[0].ID, "First message ID should match, actual messages:", messages)
	assert.Equal(t, "5", messages[1].ID, "Second message ID should match, actual messages:", messages)

	messages, err = GetMessagesAfter(roomID, 0)
	assert.NoError(t, err, "Getting messages after should not return an error.")
	assert.Equal(t, 10, len(messages), "Should return two messages the 10 messages after the specified time.")
	assert.Equal(t, "1", messages[0].ID, "First message ID should match, actual messages:", messages)
	assert.Equal(t, "2", messages[1].ID, "Second message ID should match, actual messages:", messages)

	// Make sure the thing doesn't return empty messages beyond the limit of 10
	messages, err = GetMessagesAfter(roomID, 130000)
	assert.NoError(t, err, "Getting messages after should not return an error.")
	assert.Equal(t, 3, len(messages), "Should return two messages before the specified time.")
	assert.Equal(t, "14", messages[0].ID, "First message ID should match, actual messages:", messages)
	assert.Equal(t, "15", messages[1].ID, "Second message ID should match, actual messages:", messages)

}

func TestGetMessageById(t *testing.T) {
	roomID := "room3"
	messageMap.Store(roomID, &MessageSink{
		Mutex: &sync.Mutex{},
		Messages: []Message{
			{ID: "msg1", Data: "Hello"},
			{ID: "msg2", Data: "World"},
		},
	})

	msg, err := GetMessageById(roomID, "msg1")
	assert.NoError(t, err, "Getting message by ID should not return an error")
	assert.Equal(t, "Hello", msg.Data, "Message data should match the requested message")
}

func TestDeleteMessage(t *testing.T) {
	roomID := "room1"
	messageMap.Store(roomID, &MessageSink{
		Mutex: &sync.Mutex{},
		Messages: []Message{
			{ID: "msg1", Data: "Hello"},
			{ID: "msg2", Data: "World"},
		},
	})

	err := DeleteMessage(roomID, "msg2")
	assert.NoError(t, err, "Deleting message by ID should not return an error")

	// Make sure the message has actually been deleted
	obj, _ := messageMap.Load(roomID)
	sink := obj.(*MessageSink)
	assert.Equal(t, 1, len(sink.Messages), "Message should have been deleted")
	assert.Equal(t, "msg1", sink.Messages[0].ID, "Wrong message was deleted")
}

func TestGetMessagesBeforeNoMessages(t *testing.T) {
	roomID := "room4"
	messageMap.Store(roomID, &MessageSink{
		Mutex:    &sync.Mutex{},
		Messages: []Message{},
	})

	messages, err := GetMessagesBefore(roomID, 10000)
	assert.NoError(t, err, "Getting messages before with no messages should not return an error")
	assert.Equal(t, 0, len(messages), "Should return no messages")
}

func TestGetMessagesAfterNoMessages(t *testing.T) {
	roomID := "room5"
	messageMap.Store(roomID, &MessageSink{
		Mutex:    &sync.Mutex{},
		Messages: []Message{},
	})

	messages, err := GetMessagesAfter(roomID, 10000)
	assert.NoError(t, err, "Getting messages after with no messages should not return an error")
	assert.Equal(t, 0, len(messages), "Should return no messages")
}

func TestGetMessageByIdNotFound(t *testing.T) {
	roomID := "room6"
	messageMap.Store(roomID, &MessageSink{
		Mutex:    &sync.Mutex{},
		Messages: []Message{},
	})

	_, err := GetMessageById(roomID, "nonexistent")
	assert.Error(t, err, "Getting a nonexistent message should not return an error")
}
