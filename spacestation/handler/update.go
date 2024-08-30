package handler

import (
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
)

// Action: update
func update(c *pipeshandler.Context, action struct {
	Muted    bool `json:"muted"`
	Deafened bool `json:"deafened"`
}) pipes.Event {

	connection, valid := caching.GetConnection(c.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(c, "invalid", nil)
	}

	connections, valid := caching.GetAllConnections(c.Client.Session)
	if !valid {
		return pipeshandler.ErrorResponse(c, "invalid", nil)
	}

	client := connections[c.Client.ID]
	client.ClientID = connection.ClientID
	client.Muted = action.Muted
	client.Deafened = action.Deafened
	util.Log.Println("UPDATED CLIENT", client.Data, client.ClientID, connection.ID)
	connections[c.Client.ID] = client
	caching.SaveConnections(c.Client.Session, connections)

	// Send to all
	if !SendStateUpdate(connection.ClientID, c.Client.Session, client.Muted, client.Deafened) {
		return pipeshandler.ErrorResponse(c, integration.ErrorServer, nil)
	}

	return pipeshandler.SuccessResponse(c)
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
