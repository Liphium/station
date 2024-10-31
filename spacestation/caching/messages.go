package caching

import (
	"slices"
	"sync"

	"github.com/Liphium/station/spacestation/util"
	"github.com/google/uuid"
)

type Message struct {
	ID uuid.UUID `json:"id"`

	Conversation string `json:"conversation"` // The room id
	Creation     int64  `json:"creation"`     // Unix timestamp (SET BY THE CLIENT, EXTREMELY IMPORTANT FOR SIGNATURES)
	Data         string `json:"data"`         // Encrypted data
	Edited       bool   `json:"edited"`       // Edited flag
	Sender       string `json:"sender"`       // Sender ID (of conversation token)
}

type MessageSink struct {
	Mutex    *sync.Mutex
	Messages []Message // A sorted list of messages
}

// How the messages should be sorted:
// The newest messages should have the highest indices of the message slice.
//
// Index - Creation Time - Message
// 0 - 10000 - "hello" (oldest message, at index 0)
// 1 - 20000 - "wassup"
// 2 - 30000 - "doing good?"
// 3 - 40000 - "yessir" (newest message, at highest index)

// Room id -> *MessageSink
var messageMap *sync.Map = &sync.Map{}

// Gets 10 messages before the specified time using an optimized binary search
// algorithm to make sure this doesn't take long even when there are hundreds or
// thousands of messages.
func GetMessagesBefore(room string, time int64) ([]Message, error) {

	// Get the message sink
	obj, valid := messageMap.Load(room)
	if !valid {
		return []Message{}, nil
	}
	sink := obj.(*MessageSink)

	// Copy the messages out of the sink
	sink.Mutex.Lock()
	messagesCopy := make([]Message, len(sink.Messages))
	copy(messagesCopy, sink.Messages)
	sink.Mutex.Unlock()

	// If there aren't more than 10 messages just return them as is
	if len(messagesCopy) <= 10 {

		// Remove the messages that aren't before the specified time
		return slices.DeleteFunc(messagesCopy, func(msg Message) bool {
			return msg.Creation >= time
		}), nil
	}

	// Find the nearest message to the time
	found := false
	length := len(messagesCopy)
	currentJump := length / 2
	currentIndex := currentJump - 1
	maxIndex := length - 1
	for {
		msg := messagesCopy[currentIndex]

		// Check if the current message was sent before the specified time
		if msg.Creation < time {

			// If the index is equal to zero, the process is done (no more messages to search)
			if currentIndex == 0 {
				found = false
				break
			}

			// If the index is at the top of the array, the process is done (all messages before time)
			if currentIndex == maxIndex {
				found = true
				break
			}

			// If the message sent after the current message (at an index higher than the current message)
			// is not below the time parameter, the current message is the start of the array
			if messagesCopy[currentIndex+1].Creation >= time {
				found = true
				break
			}

			// If the message was before, but not enough, jump to a higher index
			currentIndex = currentIndex + currentJump
			if currentIndex > maxIndex {
				currentIndex = maxIndex
			}
		} else {

			// If it was not sent before the specified time, jump to a lower index
			currentIndex = currentIndex - currentJump
			if currentIndex < 0 {
				currentIndex = 0
			}
		}

		// Decrease the jump by half if possible
		if currentJump <= 1 {
			currentJump = 1
		} else {
			currentJump /= 2
		}
	}

	// If no index has been found, return an empty list of messages
	if !found {
		return []Message{}, nil
	}

	// Get 10 messages before the current index
	messages := make([]Message, 10)
	copied := copy(messages, messagesCopy[currentIndex:])

	// TODO: Remove this debug messages after testing
	util.Log.Println("Copied ", copied, " messages into the array")

	return messages, nil
}
