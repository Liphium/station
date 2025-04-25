package caching

import (
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
)

// Security configuration
const MaxSharedSpaces = 10

// Conversation id -> sync.Map of SharedSpace instances
var sharedSpacesMap = &sync.Map{}

type SharedSpace struct {
	Id           string
	UnderlyingId string // Id of the Space (when pinned, so things don't get created twice)
	Name         string // Encrypted: Name of the Space
	Conversation string
	Server       string
	Mutex        *sync.Mutex
	Members      []string // Encrypted (member ids)
	Container    string   // Encrypted (Space connection container)
}

// Mutex indicates whether the mutex of the shared space should be locked
func (s *SharedSpace) ToSendable(mutex bool) SendableSharedSpace {
	if mutex {
		s.Mutex.Lock()
		defer s.Mutex.Unlock()
	}

	return SendableSharedSpace{
		UnderlyingId: s.UnderlyingId,
		Name:         s.Name,
		Members:      s.Members,
		Container:    s.Container,
	}
}

type SendableSharedSpace struct {
	UnderlyingId string   `json:"id"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	Container    string   `json:"container"`
}

// Store a new shared space, returns if it already exists and a translation for an error message (not nil = error)
func StoreSharedSpace(
	conversation string,
	server string,
	id string,
	name string,
	underlyingId string,
	container string,
) (bool, localization.Translations) {

	// Get the current map or create a new one for the conversation
	var spaceMap *sync.Map
	obj, ok := sharedSpacesMap.Load(conversation)
	if !ok {
		spaceMap = &sync.Map{}
		sharedSpacesMap.Store(conversation, spaceMap)
	} else {
		spaceMap = obj.(*sync.Map)
	}

	// Make sure there aren't too many
	count := 0
	spaceMap.Range(func(key, value any) bool {
		count++
		return true
	})
	if count >= MaxSharedSpaces {
		return false, localization.ErrorTooManySharedSpaces(MaxSharedSpaces)
	}

	// Make sure the Space isn't in there already
	if _, ok := spaceMap.Load(id); ok {
		return true, nil
	}

	// Make sure the same underlying id isn't used yet (if desired)
	if underlyingId != "-" {
		found := false
		spaceMap.Range(func(_, value any) bool {
			space := value.(*SharedSpace)
			if space.UnderlyingId == underlyingId { // Shouldn't need mutex here (the value of this never changes)
				found = true
				return false
			}
			return true
		})
		if found {
			return true, nil
		}
	}

	if !strings.HasPrefix(server, "http") {
		server = integration.Protocol + server // Add http:// or https:// in case not there
	}
	space := &SharedSpace{
		Id:           id,
		UnderlyingId: underlyingId,
		Name:         name,
		Conversation: conversation,
		Server:       server,
		Mutex:        &sync.Mutex{},
		Container:    container,
	}
	spaceMap.Store(id, space)
	startSharedSpaceInfoPuller(space)

	return false, nil
}

// Start the goroutine pulling new information about the shared space
func startSharedSpaceInfoPuller(space *SharedSpace) {
	go func() {
		for {
			// Pull new info from space station
			resp, err := integration.PostRequestNoTC(space.Server+"/info", map[string]interface{}{
				"room": space.Id,
			})
			if err != nil {
				continue
			}

			// Delete the shared space in case it is not there anymore
			if !resp["success"].(bool) {
				deleteSharedSpace(space, true)
				break
			}

			// Update in the actual shared space
			space.Mutex.Lock()
			members := resp["members"].([]string)
			if !slices.Equal(space.Members, members) {
				space.Members = members
			}
			SendEventToConversation(space.Conversation, SharedSpacesUpdateEvent(space, false))
			space.Mutex.Unlock()

			time.Sleep(time.Second * 10)
		}
	}()
}

// Create the event used to update a shared space (mutex indicates whether the mutex of the spaces should be locked)
func SharedSpacesUpdateEvent(space *SharedSpace, mutex bool) pipes.Event {
	return pipes.Event{
		Name: "shared_space",
		Data: map[string]interface{}{
			"space": space.ToSendable(mutex),
		},
	}
}

// Mutex indicates whether the mutex of the space should be locked or not
func deleteSharedSpace(space *SharedSpace, mutex bool) {
	if mutex {
		space.Mutex.Lock()
		defer space.Mutex.Unlock()
	}

	// Delete the thing
	if obj, ok := sharedSpacesMap.Load(space.Conversation); ok {
		spaceMap := obj.(*sync.Map)
		spaceMap.Delete(space.Id)
	}

	// Build the deletion event
	event := pipes.Event{
		Name: "shared_spaces_delete",
		Data: map[string]interface{}{
			"id": space.Id,
		},
	}
	if err := SendEventToConversation(space.Conversation, event); err != nil {
		util.Log.Println("ERROR: couldn't send shared space delete event:", err)
	}
}

// Rename a shared space
func RenameSharedSpace(conversation string, id string, name string) {

	// Get the space map for the conversation (or return if not there)
	obj, ok := sharedSpacesMap.Load(conversation)
	if !ok {
		return
	}
	spaceMap := obj.(*sync.Map)

	// Try getting the shared Space (or return if not there)
	spaceObj, ok := spaceMap.Load(id)
	if !ok {
		return
	}
	space := spaceObj.(*SharedSpace)

	// Rename the thing
	space.Mutex.Lock()
	defer space.Mutex.Unlock()
	space.Name = name
}
