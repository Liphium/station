package pipeshandler

import (
	"log"
	"runtime/debug"

	"github.com/Liphium/station/pipes"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
)

func NormalResponse(ctx Context, data map[string]interface{}) {
	Response(ctx, data, ctx.Instance)
}

func Response(ctx Context, data map[string]interface{}, instance *Instance) {
	err := instance.SendEventToOne(ctx.Client, pipes.Event{
		Name: "res:" + ctx.Action + ":" + ctx.ResponseId,
		Data: data,
	})
	if err != nil {
		log.Println("error while sending response to", ctx.Action, ":", err.Error())
	}
}

func SuccessResponse(ctx Context) {
	Response(ctx, map[string]interface{}{
		"success": true,
		"message": "",
	}, ctx.Instance)
}

func StatusResponse(ctx Context, status string) {
	Response(ctx, map[string]interface{}{
		"success": true,
		"message": status,
	}, ctx.Instance)
}

func ErrorResponse(ctx Context, err string) {

	if pipes.DebugLogs {
		pipeshutil.Log.Println("error with action " + ctx.Action + ": " + err)
		debug.PrintStack()
	}

	Response(ctx, map[string]interface{}{
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
