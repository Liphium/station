package pipeshandler

import (
	"github.com/Liphium/station/pipes"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
)

type Context struct {
	Client *Client                `json:"client"`
	Action string                 `json:"action"` // The action to perform
	Data   map[string]interface{} `json:"data"`
	Node   *pipes.LocalNode       `json:"-"`
}

func (instance *Instance) RegisterHandler(action string, handler func(Context)) {
	instance.routes[action] = handler
}

func (instance *Instance) Handle(ctx Context) bool {
	defer func() {
		if err := recover(); err != nil {
			ErrorResponse(ctx, "internal")
		}
	}()

	// Check if the action exists
	if instance.routes[ctx.Action] == nil {
		return false
	}

	pipeshutil.Log.Println("Handling message: " + ctx.Action)

	go instance.route(ctx)

	return true
}

func (instance *Instance) route(ctx Context) {
	defer func() {
		if err := recover(); err != nil {
			pipeshutil.Log.Println(err)
			ErrorResponse(ctx, "invalid")
		}
	}()

	instance.routes[ctx.Action](ctx)
}
