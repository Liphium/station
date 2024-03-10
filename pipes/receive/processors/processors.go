package processors

import (
	"log"

	"github.com/Liphium/station/pipes"
	"github.com/bytedance/sonic"
)

var Processors map[string]func(*pipes.Message, string) pipes.Event = make(map[string]func(*pipes.Message, string) pipes.Event)

func ProcessMarshal(message *pipes.Message, target string) []byte {
	event := ProcessEvent(message, target)

	// Marshal the event
	msg, err := sonic.Marshal(event)
	if err != nil {
		return nil
	}

	return msg
}

func ProcessEvent(message *pipes.Message, target string) pipes.Event {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error processing message: %s \n", err)
		}
	}()

	// Process the event
	if Processors[message.Event.Name] != nil {
		return Processors[message.Event.Name](message, target)
	}

	return message.Event
}
