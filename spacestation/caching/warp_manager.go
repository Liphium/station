package caching

import (
	"sync"

	"github.com/Liphium/station/pipes"
	"github.com/google/uuid"
)

// Room -> *WarpList
var warpCache *sync.Map = &sync.Map{}

type WarpList struct {
	Mutex *sync.Mutex
	List  []*WarpData // List of Warps for the room
}

type WarpData struct {
	ID        string      `json:"id"`
	Hoster    string      `json:"ho"` // Client id of the hoster
	Port      uint        `json:"p"`  // Port of the application on the hoster's device
	Mutex     *sync.Mutex `json:"-"`
	Receivers []string    `json:"-"` // Client ids of people receiving the warp
}

// Create a new Warp in a room
func NewWarp(room string, hoster string, port uint) {

	// Get the list of warps for the current room
	obj, valid := warpCache.Load(room)
	var list *WarpList
	if !valid {

		// If there isn't a list yet, create a new one
		list = &WarpList{
			Mutex: &sync.Mutex{},
			List:  []*WarpData{},
		}
		warpCache.Store(room, list)
	} else {

		// If there is one, cast obj to the list
		list = obj.(*WarpList)
	}

	// Lock the mutex to prevent concurrent reads/writes
	list.Mutex.Lock()
	defer list.Mutex.Unlock()

	// Add the warp to the list
	warp := &WarpData{
		ID:        uuid.New().String(),
		Hoster:    hoster,
		Port:      port,
		Mutex:     &sync.Mutex{},
		Receivers: []string{},
	}
	list.List = append(list.List, warp)

	SendEventToAll(room, pipes.Event{
		Name: "wp_new",
		Data: map[string]interface{}{
			"w": warp,
		},
	})
}

// Get a list of all the Warps for a specified room.
func GetWarps(room string) ([]*WarpData, error) {

	// Get the list of warps for the current room
	obj, valid := warpCache.Load(room)
	var list *WarpList
	if !valid {

		// If there isn't a list yet, create a new one
		list = &WarpList{
			Mutex: &sync.Mutex{},
			List:  []*WarpData{},
		}
		warpCache.Store(room, list)
	} else {

		// If there is one, cast obj to the list
		list = obj.(*WarpList)
	}

	// Lock the mutex to prevent concurrent reads
	list.Mutex.Lock()
	defer list.Mutex.Unlock()

	return list.List, nil
}

func InitializeWarps(client string) {

}
