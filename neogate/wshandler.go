package neogate

import (
	"github.com/Liphium/station/main/localization"
	"github.com/bytedance/sonic"
)

type Context struct {
	Client     *Client
	Action     string // The action to perform
	Locale     string // The locale of the client
	ResponseId string
	Data       []byte
	Instance   *Instance
}

// Create a handler for an action using generics (with parsing already implemented)
func CreateHandlerFor[T any](instance *Instance, action string, handler func(*Context, T) Event) {
	instance.routes[action] = func(c *Context) Event {

		// Parse the action
		var action Message[T]
		if err := sonic.Unmarshal(c.Data, &action); err != nil {
			return ErrorResponse(c, localization.ErrorInvalidRequest, err)
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

	Log.Println("Handling message: " + ctx.Action)

	go instance.route(ctx)

	return true
}

func (instance *Instance) route(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			Log.Println("recovered from error in action", ctx.Action, "by", ctx.Client.ID, ":", err)
			if err := instance.SendEventToClient(ctx.Client, ErrorResponse(ctx, localization.ErrorInvalidRequest, nil)); err != nil {
				Log.Println("couldn't send invalid event to connection after recover:", err)
			}
		}
	}()

	// Get the response from the action
	res := instance.routes[ctx.Action](ctx)

	// Send the action to the thing
	err := instance.SendEventToClient(ctx.Client, res)
	if err != nil {
		Log.Println("error while sending response to", ctx.Action, ":", err)
	}
}
