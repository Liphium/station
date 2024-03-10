package conversation

import (
	"github.com/Liphium/station/pipes"
)

func open(message *pipes.Message, target string) pipes.Event {

	/* //* This is what's happening to get the key, but for performance reasons (garbage collector) it's done later
	keys := message.Event.Data["keys"].(map[string]string)
	key := keys[fmt.Sprintf("%d", target)]
	*/

	return pipes.Event{
		Name: "conv_open:l",
		Data: map[string]interface{}{
			"success":      true,
			"conversation": message.Event.Data["conversation"],
			"members":      message.Event.Data["members"],
			"key":          message.Event.Data["keys"].(map[string]string)[target], //* Here
		},
	}
}
