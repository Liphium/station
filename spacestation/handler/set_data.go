package handler

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/spacestation/caching"
)

func SendRoomData(id string) bool {
	adapters, event, valid := GetRoomData(id, "room_data")
	if !valid {
		return false
	}

	// Send to all
	err := caching.SSNode.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.BroadcastChannel(adapters),
		Local:   true,
		Event:   event,
	})
	return err == nil
}

func GetRoomData(id string, eventName string) ([]string, pipes.Event, bool) {
	room, validRoom := caching.GetRoom(id)
	members, valid := caching.GetAllConnections(id)
	if !valid || !validRoom {
		return []string{}, pipes.Event{}, false
	}

	// Get all members
	adapters := make([]string, len(members))
	returnableMembers := make([]caching.ReturnableMember, len(members))
	i := 0
	for _, member := range members {
		returnableMembers[i] = member.ToReturnableMember()
		adapters[i] = member.Adapter
		i++
	}

	// Send to all
	return adapters, pipes.Event{
		Name: eventName,
		Data: map[string]interface{}{
			"start":   room.Start,
			"members": returnableMembers,
		},
	}, true
}
