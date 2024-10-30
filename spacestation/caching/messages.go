package caching

import (
	"slices"
	"sync"

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
	Messages []Message
}

// Room id -> *MessageSink
var messageMap *sync.Map = &sync.Map{}

func GetMessagesBefore(room string, time int64) ([]Message, error) {

	// Get the message sink
	obj, valid := messageMap.Load(room)
	if !valid {
		return []Message{}, nil
	}
	sink := obj.(*MessageSink)

	// Lock the mutex
	sink.Mutex.Lock()
	defer sink.Mutex.Unlock()

	// If there aren't more than 10 messages just return them as is
	if len(sink.Messages) <= 10 {
		return sink.Messages, nil
	}

	// Get the first message before the passed in time
	index, valid := slices.BinarySearchFunc(sink.Messages, time, func(msg Message, time int64) int {
		return int(msg.Creation - time)
	})

	return []Message{}, nil
}
