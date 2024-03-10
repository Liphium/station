package wshandler

import (
	"log"
	"runtime/debug"

	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/send"
)

func NormalResponse(message Message, data map[string]interface{}) {
	Response(message.Client.ID, message.Action, data)
}

func Response(client string, action string, data map[string]interface{}) {
	send.Client(client, pipes.Event{
		Name: action,
		Data: data,
	})
}

func SuccessResponse(message Message) {
	Response(message.Client.ID, message.Action, map[string]interface{}{
		"success": true,
		"message": "",
	})
}

func StatusResponse(message Message, status string) {
	Response(message.Client.ID, message.Action, map[string]interface{}{
		"success": true,
		"message": status,
	})
}

func ErrorResponse(message Message, err string) {

	if pipes.DebugLogs {
		log.Println("error with action " + message.Action + ": " + err)
		debug.PrintStack()
	}

	Response(message.Client.ID, message.Action, map[string]interface{}{
		"success": false,
		"message": err,
	})
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
