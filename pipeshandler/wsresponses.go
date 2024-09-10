package pipeshandler

import (
	"runtime/debug"

	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
)

func NormalResponse(ctx *Context, data map[string]interface{}) pipes.Event {
	return Response(ctx, data, ctx.Instance)
}

func SuccessResponse(ctx *Context) pipes.Event {
	return Response(ctx, map[string]interface{}{
		"success": true,
	}, ctx.Instance)
}

func ErrorResponse(ctx *Context, message localization.Translations, err error) pipes.Event {

	if pipes.DebugLogs {
		pipeshutil.Log.Println("error with action "+ctx.Action+" (", message, "): ", err)
		debug.PrintStack()
	}

	return Response(ctx, map[string]interface{}{
		"success": false,
		"message": Translate(ctx, message),
	}, ctx.Instance)
}

// Translate any message on a request
func Translate(c *Context, message localization.Translations) string {
	locale := c.Locale
	if locale == "" {
		locale = localization.DefaultLocale
	}
	msg := message[locale]
	return msg
}

func Response(ctx *Context, data map[string]interface{}, instance *Instance) pipes.Event {
	return pipes.Event{
		Name: "res:" + ctx.Action + ":" + ctx.ResponseId,
		Data: data,
	}
}
