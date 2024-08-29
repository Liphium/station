package pipeshandler

import (
	"runtime/debug"

	"github.com/Liphium/station/pipes"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
)

func NormalResponse(ctx *Context, data map[string]interface{}) pipes.Event {
	return Response(ctx, data, ctx.Instance)
}

func SuccessResponse(ctx *Context) pipes.Event {
	return Response(ctx, map[string]interface{}{
		"success": true,
		"message": "",
	}, ctx.Instance)
}

func ErrorResponse(ctx *Context, message string, err error) pipes.Event {

	if pipes.DebugLogs {
		pipeshutil.Log.Println("error with action "+ctx.Action+" (", message, "): ", err)
		debug.PrintStack()
	}

	return Response(ctx, map[string]interface{}{
		"success": false,
		"message": err,
	}, ctx.Instance)
}

func Response(ctx *Context, data map[string]interface{}, instance *Instance) pipes.Event {
	return pipes.Event{
		Name: "res:" + ctx.Action + ":" + ctx.ResponseId,
		Data: data,
	}
}
