package pipeshandler

import (
	"runtime/debug"

	"github.com/Liphium/station/pipes"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
)

func NormalResponse(ctx Context, data map[string]interface{}) {
	Response(ctx.Client, ctx.Action, data, ctx.Instance)
}

func Response(client *Client, action string, data map[string]interface{}, instance *Instance) {
	instance.SendEventToOne(client, pipes.Event{
		Name: action,
		Data: data,
	})
}

func SuccessResponse(ctx Context) {
	Response(ctx.Client, ctx.Action, map[string]interface{}{
		"success": true,
		"message": "",
	}, ctx.Instance)
}

func StatusResponse(ctx Context, status string) {
	Response(ctx.Client, ctx.Action, map[string]interface{}{
		"success": true,
		"message": status,
	}, ctx.Instance)
}

func ErrorResponse(ctx Context, err string) {

	if pipes.DebugLogs {
		pipeshutil.Log.Println("error with action " + ctx.Action + ": " + err)
		debug.PrintStack()
	}

	Response(ctx.Client, ctx.Action, map[string]interface{}{
		"success": false,
		"message": err,
	}, ctx.Instance)
}

// Returns true if one of the fields is not set
func (ctx *Context) ValidateForm(fields ...string) bool {

	for _, field := range fields {
		if ctx.Data[field] == nil {
			return true
		}
	}

	return false
}
