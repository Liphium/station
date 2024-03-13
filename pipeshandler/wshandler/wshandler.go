package wshandler

import (
	"time"

	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/send"
	"github.com/Liphium/station/pipeshandler"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
)

type Message struct {
	Client *pipeshandler.Client   `json:"client"`
	Action string                 `json:"action"` // The action to perform
	Data   map[string]interface{} `json:"data"`
}

// Routes is a map of all the routes
var Routes map[string]func(Message)

func Handle(message Message) bool {
	defer func() {
		if err := recover(); err != nil {
			ErrorResponse(message, "internal")
		}
	}()

	// Check if the action exists
	if Routes[message.Action] == nil {
		return false
	}

	pipeshutil.Log.Println("Handling message: " + message.Action)

	go Route(message.Action, message)

	return true
}

func Route(action string, message Message) {
	defer func() {
		if err := recover(); err != nil {
			pipeshutil.Log.Println(err)
			ErrorResponse(message, "invalid")
		}
	}()

	Routes[message.Action](message)
}

func Initialize() {
	Routes = make(map[string]func(Message))
}

func TestConnection() {
	go func() {
		for {
			time.Sleep(time.Second * 5)

			// Send ping
			send.Pipe(send.ProtocolWS, pipes.Message{
				Channel: pipes.BroadcastChannel([]string{"1", "3"}),
				Event: pipes.Event{
					Name: "ping",
					Data: map[string]interface{}{
						"node": pipes.CurrentNode.ID,
					},
				},
			})
		}
	}()
}
