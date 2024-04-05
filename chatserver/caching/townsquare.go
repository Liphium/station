package caching

import (
	"sync"

	"github.com/Liphium/station/pipes"
)

type TownsquareMember struct {
	Mutex    *sync.Mutex
	Account  string
	Username string
	Viewing  bool
}

// ! Always use cost 1
var townsquareCache = &sync.Map{} // Account ID -> Townsquare status
var townsquareMessageCache = &sync.Map{}

// Add someone to townsquare
func JoinTownsquare(id string, username string) {

	if _, ok := townsquareCache.Load(id); ok {
		return
	}

	townsquareCache.Store(id, &TownsquareMember{
		Mutex:    &sync.Mutex{},
		Account:  id,
		Username: username,
		Viewing:  false,
	})

	// Tell everyone about the join
	SendTownsquareEvent(pipes.Event{
		Name: "townsquare_join",
		Data: map[string]interface{}{
			"id":       id,
			"username": username,
		},
	})
}

// Remove someone from townsquare
func LeaveTownsquare(id string) {

	if _, ok := townsquareCache.Load(id); !ok {
		return
	}
	townsquareCache.Delete(id)

	// Send leave event to everyone
	SendTownsquareEvent(pipes.Event{
		Name: "townsquare_leave",
		Data: map[string]interface{}{
			"id": id,
		},
	})
}

// Toggle the viewing state of someone in townsquare
func SetTownsquareViewing(id string, state bool) {

	obj, ok := townsquareCache.Load(id)
	if !ok {
		return
	}
	member := obj.(*TownsquareMember)

	// Make sure this doesn't happen concurrently
	member.Mutex.Lock()
	defer member.Mutex.Unlock()

	// Change state and notify everyone
	member.Viewing = state
	if member.Viewing {
		SendTownsquareEvent(pipes.Event{
			Name: "townsquare_open",
			Data: map[string]interface{}{
				"id": id,
			},
		})
	} else {
		SendTownsquareEvent(pipes.Event{
			Name: "townsquare_close",
			Data: map[string]interface{}{
				"id": id,
			},
		})
	}
}

// Send an event to all people in townsquare
func SendTownsquareEvent(event pipes.Event) {

	// Iterate through all members and send the event to the client
	townsquareCache.Range(func(key, value any) bool {
		CSNode.SendClient(key.(string), event)
		return true
	})
}
