package pipeshandler

import (
	"runtime/debug"

	"github.com/Liphium/station/pipes"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
)

func NormalResponse(ctx Context, data map[string]interface{}) {
	Response(ctx.Client.ID, ctx.Action, data, ctx.Node)
}

func Response(client string, action string, data map[string]interface{}, local *pipes.LocalNode) {
	local.SendClient(client, pipes.Event{
		Name: action,
		Data: data,
	})
}

func SuccessResponse(ctx Context) {
	Response(ctx.Client.ID, ctx.Action, map[string]interface{}{
		"success": true,
		"message": "",
	}, ctx.Node)
}

func StatusResponse(ctx Context, status string) {
	Response(ctx.Client.ID, ctx.Action, map[string]interface{}{
		"success": true,
		"message": status,
	}, ctx.Node)
}

func ErrorResponse(ctx Context, err string) {

	if pipes.DebugLogs {
		pipeshutil.Log.Println("error with action " + ctx.Action + ": " + err)
		debug.PrintStack()
	}

	Response(ctx.Client.ID, ctx.Action, map[string]interface{}{
		"success": false,
		"message": err,
	}, ctx.Node)
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
