package caching

import (
	"errors"
	"slices"
	"sync"

	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/google/uuid"
)

// Room -> *WarpList
var warpCache *sync.Map = &sync.Map{}

type WarpList struct {
	List *sync.Map // Map of all Warps: ID -> *WarpData
}

type WarpData struct {
	ID        string      `json:"id"`
	Hoster    string      `json:"ho"` // Client id of the hoster
	Port      uint        `json:"p"`  // Port of the application on the hoster's device
	Mutex     *sync.Mutex `json:"-"`
	Receivers []string    `json:"-"` // Client ids of people receiving the warp
}

// Create a new Warp in a room
func NewWarp(room string, hoster string, port uint) (string, error) {

	// Get the list of warps for the current room
	obj, valid := warpCache.Load(room)
	var list *WarpList
	if !valid {

		// If there isn't a list yet, create a new one
		list = &WarpList{
			List: &sync.Map{},
		}
		warpCache.Store(room, list)
	} else {

		// If there is one, cast obj to the list
		list = obj.(*WarpList)
	}

	// Add the warp to the list
	warp := &WarpData{
		ID:        uuid.New().String(),
		Hoster:    hoster,
		Port:      port,
		Mutex:     &sync.Mutex{},
		Receivers: []string{},
	}
	list.List.Store(warp.ID, warp)

	return warp.ID, SendEventToAll(room, pipes.Event{
		Name: "wp_new",
		Data: map[string]interface{}{
			"w": warp.ID,
			"h": hoster,
			"p": port,
		},
	})
}

// Get a list of all the Warps for a specified room.
func RangeOverWarps(room string, rangeFunc func(warpId string, w *WarpData) bool) error {

	// Get the list of warps for the current room
	obj, valid := warpCache.Load(room)
	var list *WarpList
	if !valid {

		// If there isn't a list yet, create a new one
		list = &WarpList{
			List: &sync.Map{},
		}
		warpCache.Store(room, list)
	} else {

		// If there is one, cast obj to the list
		list = obj.(*WarpList)
	}

	// Range over everything
	list.List.Range(func(key, value any) bool {
		return rangeFunc(key.(string), value.(*WarpData))
	})

	return nil
}

// Send a client all the warps upon joining a Space.
func InitializeWarps(client *pipeshandler.Client) {
	RangeOverWarps(client.Session, func(warpId string, w *WarpData) bool {
		return SSNode.SendClient(client.ID, pipes.Event{
			Name: "wp_new",
			Data: map[string]interface{}{
				"w": w.ID,
				"h": w.Hoster,
				"p": w.Port,
			},
		}) != nil
	})
}

// Stop a warp and disconnect all clients from it.
func StopWarp(room string, warp string) error {

	// Get the list of warps for the current room
	obj, valid := warpCache.Load(room)
	if !valid {
		return errors.New("no warps found")
	}
	list := obj.(*WarpList)

	// Add the warp to the list
	list.List.Delete(warp)

	return SendEventToAll(room, pipes.Event{
		Name: "wp_end",
		Data: map[string]interface{}{
			"w": warp,
		},
	})
}

// Stop all warps in a room created by a certain client id.
func StopWarpsBy(room string, client string) {
	RangeOverWarps(room, func(warpId string, w *WarpData) bool {
		if w.Hoster == client {
			StopWarp(room, warpId)
		}
		return true
	})
}

// Get any warp in a room by ID.
func GetWarp(room string, warp string) (*WarpData, error) {

	// Get the list of warps for the current room
	obj, valid := warpCache.Load(room)
	if !valid {
		return nil, errors.New("no warps found")
	}
	list := obj.(*WarpList)

	// Get the warp by id
	w, valid := list.List.Load(warp)
	if !valid {
		return nil, errors.New("warp not found")
	}
	return w.(*WarpData), nil
}

// Remove a client from receiving packets through Warp.
func RemoveClientFromWarp(id string, roomId string, warpId string) error {

	// Get the warp
	warp, err := GetWarp(roomId, warpId)
	if err != nil {
		return err
	}

	// Make sure there are no concurrent reads/writes
	warp.Mutex.Lock()
	defer warp.Mutex.Unlock()

	// Remove the receiver from the list of receivers
	warp.Receivers = slices.DeleteFunc(warp.Receivers, func(e string) bool {
		return e == id
	})

	// Tell the hoster the client has disconnected
	return SSNode.SendClient(warp.Hoster, pipes.Event{
		Name: "wp_disconnected",
		Data: map[string]interface{}{
			"w": warpId,
			"c": id,
		},
	})
}

// Send an event to all receivers of a warp.
func SendEventToWarp(room string, warpId string, event pipes.Event) error {

	// Get warp
	warp, err := GetWarp(room, warpId)
	if err != nil {
		return err
	}

	// Send event to all receivers of the warp using pipes
	return SSNode.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.BroadcastChannel(warp.Receivers),
		Local:   true,
		Event:   event,
	})
}
