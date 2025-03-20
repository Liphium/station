package caching

import (
	"errors"

	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/spacestation/util"
	"github.com/dgraph-io/ristretto"
)

type RoomConnection struct {
	Connected bool   `json:"-"`
	Adapter   string `json:"id"`   // Also the client id
	Data      string `json:"data"` // The account id of the client (encrypted)
	Signature string `json:"sign"` // Client id + Account id signed with private key of the client (to proof the account id is correct)

	StudioConnection bool `json:"st"`   // Whether or not the client is connected to Studio
	Muted            bool `json:"mute"` // If the client is muted
	Deafened         bool `json:"deaf"` // If the client is deafened or not
}

// Member (Connection) ID -> Connections
type RoomConnections map[string]RoomConnection

// ! For setting please ALWAYS use cost 1
var roomConnectionsCache *ristretto.Cache

func setupRoomConnectionsCache() {
	var err error
	roomConnectionsCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e4,     // expecting to store 1k room connections
		MaxCost:     1 << 30, // maximum cost of cache is 1GB
		BufferItems: 64,      // Some random number, check docs
		OnEvict: func(item *ristretto.Item) {
			util.Log.Println("[cache] room", item.Key, "'s connections were deleted")
		},
	})

	if err != nil {
		panic(err)
	}
}

// Sets the member data
func SetMemberData(roomID string, connectionId string, data string, signature string) bool {

	// Get the room the member is a part of
	room, valid := GetRoom(roomID)
	if !valid {
		return false
	}
	room.Mutex.Lock() // Make sure the map is not modified at the same time
	defer room.Mutex.Unlock()

	// Update the connection
	obj, valid := roomConnectionsCache.Get(roomID)
	if !valid {
		return false
	}
	connections := obj.(RoomConnections)
	if connections[connectionId].Connected {
		return false
	}
	connections[connectionId] = RoomConnection{
		Connected:        true,
		Adapter:          connectionId,
		Data:             data,
		Signature:        signature,
		StudioConnection: false,
		Muted:            true,
		Deafened:         false,
	}

	// Update the connections accordingly
	roomConnectionsCache.Set(roomID, connections, 1)
	roomConnectionsCache.Wait()

	// Send a room data update to everyone
	return SendRoomUpdateToAll(room, connections)
}

// UpdateMemberData updates the member data, specifically for studio-related states (only set the values you want to update, rest can be nil)
func UpdateMemberData(roomID string, connectionId string, studioConnection *bool, muted *bool, deafened *bool) bool {
	// Get the room the member is a part of
	room, valid := GetRoom(roomID)
	if !valid {
		return false
	}
	room.Mutex.Lock() // Make sure the map is not modified at the same time
	defer room.Mutex.Unlock()

	// Update the connection
	obj, valid := roomConnectionsCache.Get(roomID)
	if !valid {
		return false
	}
	connections := obj.(RoomConnections)
	connection, exists := connections[connectionId]
	if !exists {
		return false
	}

	// Update the studio-related states
	if studioConnection != nil {
		connection.StudioConnection = *studioConnection
	}
	if muted != nil {
		connection.Muted = *muted
	}
	if deafened != nil {
		connection.Deafened = *deafened
	}
	connections[connectionId] = connection

	// Update the connections accordingly
	roomConnectionsCache.Set(roomID, connections, 1)
	roomConnectionsCache.Wait()

	// Send a room data update to everyone
	return SendRoomUpdateToAll(room, connections)
}

// Remove a member from a room
func RemoveMember(roomID string, connectionId string) bool {

	// Get the room where the member should be removed
	room, valid := GetRoom(roomID)
	if !valid {
		return false
	}
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	// Delete the member
	obj, valid := roomConnectionsCache.Get(roomID)
	if !valid {
		return false
	}
	connections := obj.(RoomConnections)
	delete(connections, connectionId)
	if len(connections) == 0 {
		DeleteRoom(roomID)
		return true
	}

	// Update the connections map
	roomConnectionsCache.Set(roomID, connections, 1)
	roomConnectionsCache.Wait()

	// Send a room update event to everyone
	return SendRoomUpdateToAll(room, connections)
}

// Get all connections from a room
func GetAllConnections(room string) (RoomConnections, bool) {

	connections, found := roomConnectionsCache.Get(room)

	if !found {
		return nil, false
	}

	return connections.(RoomConnections), true
}

// Get all adapters from a room
func GetAllAdapters(room string) ([]string, bool) {

	connections, valid := GetAllConnections(room)
	if !valid {
		return nil, false
	}

	adapters := make([]string, len(connections))
	i := 0
	for key := range connections {
		adapters[i] = key
		i++
	}

	return adapters, true
}

// Save changes in a room
func SaveConnections(roomId string, connections RoomConnections) bool {

	room, valid := GetRoom(roomId)
	if !valid {
		return false
	}
	room.Mutex.Lock()

	room, valid = GetRoom(roomId)
	if !valid {
		room.Mutex.Unlock()
		return false
	}

	// Refresh room
	roomConnectionsCache.Set(roomId, connections, 1)
	roomConnectionsCache.Wait()
	room.Mutex.Unlock()

	return true
}

// Get the event for a room data update (returns adapter, event and if there was an error)
func GetRoomDataEvent(room Room, members RoomConnections) ([]string, pipes.Event, bool) {

	// Get all the adapters and members
	connections := make([]RoomConnection, len(members))
	adapters := make([]string, len(members))
	i := 0
	for _, member := range members {
		connections[i] = member
		adapters[i] = member.Adapter
		i++
	}

	// Return the event
	return adapters, pipes.Event{
		Name: "room_data",
		Data: map[string]interface{}{
			"start":   room.Start,
			"members": connections,
		},
	}, true
}

// Send a room update event to all clients in a room
func SendRoomUpdateToAll(room Room, connections RoomConnections) bool {

	// Get the actual event and adapters
	adapters, event, valid := GetRoomDataEvent(room, connections)
	if !valid {
		return false
	}

	// Send the room update event using pipes
	err := SSNode.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.BroadcastChannel(adapters),
		Local:   true,
		Event:   event,
	})
	return err == nil
}

// Send an event to all members of a room.
func SendEventToAll(room string, event pipes.Event) error {

	// Get all adapters for the people in the room
	adapters, valid := GetAllAdapters(room)
	if !valid {
		return errors.New("adapters couldn't be found for this room")
	}

	// Send the actual event using pipes
	if err := SSNode.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.BroadcastChannel(adapters),
		Local:   true,
		Event:   event,
	}); err != nil {
		return err
	}

	return nil
}

// Make a value a pointer (helper function)
func Ptr[T any](val T) *T {
	return &val
}
