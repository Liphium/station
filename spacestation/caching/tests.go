package caching

import (
	"testing"
	"time"

	"github.com/Liphium/station/spacestation/util"
)

// TestRooms tests the caching of rooms
func TestRooms(t *testing.T) {

	if roomsCache == nil {
		t.Error("Rooms cache not initialized")
		return
	}

	// Test caching
	CreateRoom("id")

	room, valid := GetRoom("id")
	if !valid {
		t.Error("Room not found")
		return
	} else {
		if room.ID != "id" {
			t.Error("Room has wrong ID")
		}
	}

	for i := 0; i < 10; i++ {
		go func() {
			valid := JoinRoom("id", util.GenerateToken(5))
			if !valid {
				t.Error("Room not found")
			}
		}()
	}

	time.Sleep(time.Millisecond * 500)
	connections, valid := GetAllConnections("id")
	if !valid {
		t.Error("Connections couldn't be retrieved")
	} else {
		if len(connections) != 10 {
			t.Errorf("Room has wrong number of members (expected 10, got %d)", len(connections))
		}
	}
}
