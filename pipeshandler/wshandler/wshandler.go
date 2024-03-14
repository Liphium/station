package wshandler

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
)

type Message struct {
	Client *pipeshandler.Client   `json:"client"`
	Action string                 `json:"action"` // The action to perform
	Data   map[string]interface{} `json:"data"`
	Node   *pipes.LocalNode       `json:"-"`
}

// Routes is a map of all the routes
var Routes map[string]map[string]func(Message)

func RegisterHandler(local *pipes.LocalNode, action string, handler func(Message)) {
	if Routes[local.ID] == nil {
		Routes[local.ID] = make(map[string]func(Message))
	}

	Routes[local.ID][action] = handler
}

func Handle(message Message) bool {
	defer func() {
		if err := recover(); err != nil {
			ErrorResponse(message, "internal")
		}
	}()

	// Check if the action exists
	if Routes[message.Node.ID] == nil || Routes[message.Node.ID][message.Action] == nil {
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

	Routes[message.Node.ID][message.Action](message)
}

func Initialize() {
	Routes = make(map[string]map[string]func(Message))
}
