package wshandler

import (
	"runtime/debug"

	"github.com/Liphium/station/pipes"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
)

func NormalResponse(message Message, data map[string]interface{}) {
	Response(message.Client.ID, message.Action, data, message.Node)
}

func Response(client string, action string, data map[string]interface{}, local *pipes.LocalNode) {
	local.SendClient(client, pipes.Event{
		Name: action,
		Data: data,
	})
}

func SuccessResponse(message Message) {
	Response(message.Client.ID, message.Action, map[string]interface{}{
		"success": true,
		"message": "",
	}, message.Node)
}

func StatusResponse(message Message, status string) {
	Response(message.Client.ID, message.Action, map[string]interface{}{
		"success": true,
		"message": status,
	}, message.Node)
}

func ErrorResponse(message Message, err string) {

	if pipes.DebugLogs {
		pipeshutil.Log.Println("error with action " + message.Action + ": " + err)
		debug.PrintStack()
	}

	Response(message.Client.ID, message.Action, map[string]interface{}{
		"success": false,
		"message": err,
	}, message.Node)
}

// Returns true if one of the fields is not set
func (message *Message) ValidateForm(fields ...string) bool {

	for _, field := range fields {
		if message.Data[field] == nil {
			return true
		}
	}

	return false
}
