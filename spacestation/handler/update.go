package handler

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler/wshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
)

// Action: update
func update(message wshandler.Message) {

	if message.ValidateForm("muted", "deafened") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	connection, valid := caching.GetConnection(message.Client.ID)
	if !valid {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	connections, valid := caching.GetAllConnections(message.Client.Session)
	if !valid {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	client := connections[message.Client.ID]
	client.ClientID = connection.ClientID
	client.Muted = message.Data["muted"].(bool)
	client.Deafened = message.Data["deafened"].(bool)
	util.Log.Println("UPDATED CLIENT", client.Data, client.ClientID, connection.ID)
	connections[message.Client.ID] = client
	caching.SaveConnections(message.Client.Session, connections)

	// Send to all
	if !SendStateUpdate(connection.ClientID, message.Client.Session, client.Muted, client.Deafened) {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	wshandler.SuccessResponse(message)
}

func SendStateUpdate(member string, room string, muted bool, deafened bool) bool {

	// Get all adapters
	adapters, valid := caching.GetAllAdapters(room)
	if !valid {
		return false
	}

	// Send to all
	err := caching.Node.Pipe(pipes.ProtocolWS, pipes.Message{
		Local:   true,
		Channel: pipes.BroadcastChannel(adapters),
		Event: pipes.Event{
			Name: "member_update",
			Data: map[string]interface{}{
				"member":   member,
				"muted":    muted,
				"deafened": deafened,
			},
		},
	})
	return err == nil
}
