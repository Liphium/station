package handler

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
)

// Action: update
func update(ctx pipeshandler.Context) {

	if ctx.ValidateForm("muted", "deafened") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	connection, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	connections, valid := caching.GetAllConnections(ctx.Client.Session)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	client := connections[ctx.Client.ID]
	client.ClientID = connection.ClientID
	client.Muted = ctx.Data["muted"].(bool)
	client.Deafened = ctx.Data["deafened"].(bool)
	util.Log.Println("UPDATED CLIENT", client.Data, client.ClientID, connection.ID)
	connections[ctx.Client.ID] = client
	caching.SaveConnections(ctx.Client.Session, connections)

	// Send to all
	if !SendStateUpdate(connection.ClientID, ctx.Client.Session, client.Muted, client.Deafened) {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	pipeshandler.SuccessResponse(ctx)
}

func SendStateUpdate(member string, room string, muted bool, deafened bool) bool {

	// Get all adapters
	adapters, valid := caching.GetAllAdapters(room)
	if !valid {
		return false
	}

	// Send to all
	err := caching.SSNode.Pipe(pipes.ProtocolWS, pipes.Message{
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
