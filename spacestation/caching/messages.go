package caching

import (
	"slices"
	"sync"
	"unsafe"

	"github.com/Liphium/station/spacestation/util"
)

type Message struct {
	ID string `json:"id"`

	Conversation string `json:"conversation"` // The room id
	Creation     int64  `json:"creation"`     // Unix timestamp (SET BY THE CLIENT, EXTREMELY IMPORTANT FOR SIGNATURES)
	Data         string `json:"data"`         // Encrypted data
	Edited       bool   `json:"edited"`       // Edited flag
	Sender       string `json:"sender"`       // Sender ID (of conversation token)
}

func CheckSize(message string) bool {
	return unsafe.Sizeof(message) > 1000*6
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
	util.Log.Println("Copied ", copied, " messages into the array (before)")

	return messages, nil
}

// Gets 10 messages after the specified time using an optimized binary search
// algorithm to make sure this doesn't take long even when there are hundreds or
// thousands of messages.
func GetMessagesAfter(room string, time int64) ([]Message, error) {

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
			return msg.Creation <= time
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

			// If the index is equal to zero, the process is done (no messages after time)
			if currentIndex == 0 {
				found = true
				break
			}

			// If the index is at the top of the array, the process is done (no more left to search)
			if currentIndex == maxIndex {
				found = false
				break
			}

			// If the message sent after the current message (at an index higher than the current message)
			// is not below the time parameter, the current message is the start of the array
			if messagesCopy[currentIndex+1].Creation < time {
				found = true
				break
			}

			// If the message was before, but not enough, jump to a higher index
			currentIndex = currentIndex - currentJump
			if currentIndex < 0 {
				currentIndex = 0
			}
		} else {

			// If it was not sent before the specified time, jump to a lower index
			currentIndex = currentIndex + currentJump
			if currentIndex > maxIndex {
				currentIndex = maxIndex
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
	copied := copy(messages, messagesCopy[:currentIndex])

	// TODO: Remove this debug messages after testing
	util.Log.Println("Copied ", copied, " messages into the array (after)")

	return messages, nil
}

// Get a message by its id. This function can be a little slow as it iterates
// through all the messages in the chat.
func GetMessageById(room string, id string) (Message, error) {

	// Get the message sink
	obj, valid := messageMap.Load(room)
	if !valid {
		return Message{}, nil
	}
	sink := obj.(*MessageSink)

	// Copy the messages out of the sink (to unlock the mutex faster)
	sink.Mutex.Lock()
	messagesCopy := make([]Message, len(sink.Messages))
	copy(messagesCopy, sink.Messages)
	sink.Mutex.Unlock()

	// Find the current message
	index := slices.IndexFunc(messagesCopy, func(msg Message) bool {
		return msg.ID == id
	})

	return messagesCopy[index], nil
}

// Add message to the room.
func AddMessage(room string, msg Message) error {

	// Get the message sink
	obj, valid := messageMap.Load(room)
	if !valid {
		return nil
	}
	sink := obj.(*MessageSink)

	// Lock the mutex and make sure
	sink.Mutex.Lock()
	defer sink.Mutex.Unlock()

	return nil
}
