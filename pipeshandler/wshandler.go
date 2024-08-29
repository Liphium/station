package pipeshandler

import (
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipes"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
	"github.com/bytedance/sonic"
)

type Context struct {
	Client     *Client
	Action     string // The action to perform
	ResponseId string
	Data       []byte
	Node       *pipes.LocalNode
	Instance   *Instance
}

// Create a handler for an action using generics (with parsing already implemented)
func CreateHandlerFor[T any](instance *Instance, action string, handler func(*Context, T) pipes.Event) {
	instance.routes[action] = func(c *Context) pipes.Event {

		// Parse the action
		var action Message[T]
		if err := sonic.Unmarshal(c.Data, &action); err != nil {
			return ErrorResponse(c, localization.InvalidRequest, err)
		}

		// Let the handler handle it (literally)
		return handler(c, action.Data)
	}
}

func (instance *Instance) Handle(ctx *Context) bool {

	// Check if the action exists
	if instance.routes[ctx.Action] == nil {
		return false
	}

	pipeshutil.Log.Println("Handling message: " + ctx.Action)

	go instance.route(ctx)

	return true
}

func (instance *Instance) route(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			pipeshutil.Log.Println("recovered from error in action", ctx.Action, "by", ctx.Client.ID, ":", err)
			if err := instance.SendEventToOne(ctx.Client, ErrorResponse(ctx, localization.InvalidRequest, nil)); err != nil {
				pipeshutil.Log.Println("couldn't send invalid event to connection after recover:", err)
			}
		}
	}()

	// Get the response from the action
	res := instance.routes[ctx.Action](ctx)

	// Send the action to the thing
	err := instance.SendEventToOne(ctx.Client, res)
	if err != nil {
		pipeshutil.Log.Println("error while sending response to", ctx.Action, ":", err)
	}
}
