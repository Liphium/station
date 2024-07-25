package caching

import (
	"context"
	"net"

	"github.com/Liphium/station/spacestation/util"
	"github.com/dgraph-io/ristretto"
	"github.com/livekit/protocol/livekit"
)

type RoomConnection struct {
	Connected      bool
	Connection     *net.UDPAddr
	Adapter        string
	Key            *[]byte
	ClientID       string
	CurrentSession string
	Data           string

	//* Client status
	Muted    bool
	Deafened bool
}

func (r *RoomConnection) ToReturnableMember() ReturnableMember {
	return ReturnableMember{
		ID:       r.Data + ":" + r.ClientID,
		Muted:    r.Muted,
		Deafened: r.Deafened,
	}
}

// TODO: Implement as standard
type ReturnableMember struct {
	ID       string `json:"id"` // Syntax: data:clientID
	Muted    bool   `json:"muted"`
	Deafened bool   `json:"deafened"`
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

// JoinRoom adds a member to a room in the cache
func EnterUDP(roomID string, connectionId string, clientId string, addr *net.UDPAddr, key *[]byte) bool {

	room, valid := GetRoom(roomID)
	if !valid {
		return false
	}
	room.Mutex.Lock()

	room, valid = GetRoom(roomID)
	if !valid {
		room.Mutex.Unlock()
		return false
	}

	obj, valid := roomConnectionsCache.Get(roomID)
	if !valid {
		room.Mutex.Unlock()
		return false
	}
	connections := obj.(RoomConnections)
	conn := connections[connectionId]
	if conn.Connected {
		util.Log.Println("Error: Connection already exists")
		room.Mutex.Unlock()
		return false
	}
	connections[connectionId] = RoomConnection{
		Connected:      true,
		Connection:     addr,
		ClientID:       clientId,
		Data:           conn.Data,
		CurrentSession: "",
		Adapter:        connectionId,
		Key:            key,
	}

	// Refresh room
	roomConnectionsCache.Set(roomID, connections, 1)
	roomConnectionsCache.Wait()
	room.Mutex.Unlock()

	return true
}

// Sets the member data
func SetMemberData(roomID string, connectionId string, clientId string, data string) bool {

	room, valid := GetRoom(roomID)
	if !valid {
		return false
	}
	room.Mutex.Lock()

	room, valid = GetRoom(roomID)
	if !valid {
		room.Mutex.Unlock()
		return false
	}

	obj, valid := roomConnectionsCache.Get(roomID)
	if !valid {
		room.Mutex.Unlock()
		return false
	}
	connections := obj.(RoomConnections)
	if connections[connectionId].Connected {
		room.Mutex.Unlock()
		return false
	}
	connections[connectionId] = RoomConnection{
		Connected:      false,
		Connection:     nil,
		Adapter:        connectionId,
		CurrentSession: "",
		ClientID:       clientId,
		Data:           data,
	}

	// Refresh room
	roomConnectionsCache.Set(roomID, connections, 1)
	roomConnectionsCache.Wait()
	room.Mutex.Unlock()

	return true
}

func RemoveMember(roomID string, connectionId string) bool {

	room, valid := GetRoom(roomID)
	if !valid {
		return false
	}
	room.Mutex.Lock()

	room, valid = GetRoom(roomID)
	if !valid {
		room.Mutex.Unlock()
		return false
	}

	obj, valid := roomConnectionsCache.Get(roomID)
	if !valid {
		room.Mutex.Unlock()
		return false
	}
	connections := obj.(RoomConnections)
	delete(connections, connectionId)

	if len(connections) == 0 {
		DeleteRoom(roomID)

		// Also delete livekit room
		_, err := RoomClient.DeleteRoom(context.Background(), &livekit.DeleteRoomRequest{
			Room: room.ID,
		})
		if err != nil {
			util.Log.Println("couldn't close livekit room:", err)
		}

		return true
	}

	// Refresh room
	roomConnectionsCache.Set(roomID, connections, 1)
	roomConnectionsCache.Wait()
	room.Mutex.Unlock()

	// Get connection (for removing from the livekit room)
	conn, valid := GetConnection(connectionId)
	if !valid {
		util.Log.Println("couldn't get connection during room deletion")
		return true
	}

	// Remove the member from the livekit room
	_, err := RoomClient.RemoveParticipant(context.Background(), &livekit.RoomParticipantIdentity{
		Room:     room.ID,
		Identity: conn.ClientID,
	})
	if err != nil {
		util.Log.Println("couldn't remove member from livekit room:", err)
	}

	return true
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
