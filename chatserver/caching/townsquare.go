package caching

import (
	"sync"

	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/pipes"
)

type TownsquareMember struct {
	Mutex     *sync.Mutex
	Account   string
	PublicKey string // For signature verification
	Username  string
	Viewing   bool
}

// ! Always use cost 1
var townsquareCache = &sync.Map{} // Account ID -> Townsquare status

// Message system
type TownsquareMessage struct {
	ID          int64  `json:"i"`
	Sender      string `json:"s"`
	Content     string `json:"c"`
	Attachments string `json:"a"`
	Timestamp   int64  `json:"t"`
}

var counterMutex = &sync.Mutex{}
var messageCounter int64 = 0
var townsquareMessageCache = &sync.Map{} // Message ID -> Message

// Add someone to townsquare
func JoinTownsquare(id string, username string, key string) {

	if _, ok := townsquareCache.Load(id); ok {
		return
	}

	townsquareCache.Store(id, &TownsquareMember{
		Mutex:     &sync.Mutex{},
		Account:   id,
		PublicKey: key,
		Username:  username,
		Viewing:   false,
	})

	// Tell everyone about the join
	SendTownsquareEvent(townsquareJoinEvent(id, username, key))
}

// Remove someone from townsquare
func LeaveTownsquare(id string) {

	if _, ok := townsquareCache.Load(id); !ok {
		return
	}
	townsquareCache.Delete(id)

	// Send leave event to everyone
	SendTownsquareEvent(pipes.Event{
		Name: "ts_member_leave",
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
		SendTownsquareEvent(townsquareOpenEvent(id))
	} else {
		SendTownsquareEvent(pipes.Event{
			Name: "ts_member_close",
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
		err := CSNode.SendClient(key.(string), event)
		if err != nil {
			util.Log.Println("error while sending event", event, err)
		}
		return true
	})
}

func TownsquareMessageId() int64 {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	messageCounter++
	return messageCounter
}

// Save a message to the cache
func SaveTownsquareMessage(id int64, message TownsquareMessage) {
	townsquareMessageCache.Store(id, message)
}

// Get a townsquare message from the cache
func GetTownsquareMessage(id int64) (TownsquareMessage, bool) {
	obj, ok := townsquareMessageCache.Load(id)
	if !ok {
		return TownsquareMessage{}, ok
	}
	return obj.(TownsquareMessage), ok
}

// Send an event to all people in townsquare
func SendTownsquareMessageEvent(event pipes.Event) {

	// Iterate through all members and send the event to the client
	townsquareCache.Range(func(key, value any) bool {
		member := value.(*TownsquareMember)
		member.Mutex.Lock()
		defer member.Mutex.Unlock()

		if member.Viewing {
			CSNode.SendClient(key.(string), event)
		}
		return true
	})
}

// Tell a client about all people in townsquare
func SendAllTownsquareMembers(account string) {
	townsquareCache.Range(func(key, value any) bool {
		member := value.(*TownsquareMember)
		err := CSNode.SendClient(account, townsquareJoinEvent(member.Account, member.Username, member.PublicKey))
		if err != nil {
			util.Log.Println("error while sending townsquare member join", member.Account, ":", err)
		}

		member.Mutex.Lock()
		defer member.Mutex.Unlock()

		if member.Viewing {
			err = CSNode.SendClient(account, townsquareOpenEvent(member.Account))
			if err != nil {
				util.Log.Println("error while sending townsquare member open", member.Account, ":", err)
			}
		}

		return true
	})
}

// Send a client the latest messages
func SendLatestMessages(account string) {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	SendMessages(account, messageCounter)
}

// Send a client messages with an offset
func SendMessages(account string, before int64) {

}

func townsquareOpenEvent(id string) pipes.Event {
	return pipes.Event{
		Name: "ts_member_open",
		Data: map[string]interface{}{
			"id": id,
		},
	}
}

func townsquareJoinEvent(id string, name string, key string) pipes.Event {
	return pipes.Event{
		Name: "ts_member_join",
		Data: map[string]interface{}{
			"id":   id,
			"name": name,
			"key":  key,
		},
	}
}
