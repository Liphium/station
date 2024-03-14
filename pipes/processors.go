package pipes

import (
	"github.com/bytedance/sonic"
)

func (local *LocalNode) ProcessMarshal(message *Message, target string) []byte {
	event := local.ProcessEvent(message, target)

	// Marshal the event
	msg, err := sonic.Marshal(event)
	if err != nil {
		return nil
	}

	return msg
}

func (local *LocalNode) ProcessEvent(message *Message, target string) Event {
	defer func() {
		if err := recover(); err != nil {
			Log.Printf("Error processing message: %s \n", err)
		}
	}()

	// Process the event
	if local.Processors[message.Event.Name] != nil {
		return local.Processors[message.Event.Name](message, target)
	}

	return message.Event
}
