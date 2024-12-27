package caching

import (
	"slices"
	"sync"

	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/spacestation/util"
)

// ! For setting please ALWAYS use cost 1
// Room ID -> Table
var tablesCache *sync.Map = &sync.Map{}

type TableData struct {
	Mutex         *sync.Mutex
	Room          string
	MemberCount   int
	Members       *sync.Map    // Client ID -> Client info
	Objects       *sync.Map    // Cache for all objects on the table (Object ID -> Object)
	highestObject *TableObject // Object with the highest layer (for swapping and creating)
}

type TableMember struct {
	Client         string  // Client ID
	Color          float64 // Color of their cursor
	SelectedObject string  // The id of the currently selected object
	Enabled        bool    // If events should currently be sent to the member
}

// * Table management
func JoinTable(room string, client string, color float64) localization.Translations {

	obj, valid := tablesCache.Load(room)
	var table *TableData
	if !valid {

		// Create table
		table = &TableData{
			Mutex:       &sync.Mutex{},
			Room:        room,
			Members:     &sync.Map{},
			Objects:     &sync.Map{},
			MemberCount: 0,
		}
		tablesCache.Store(room, table)
	} else {
		table = obj.(*TableData)
	}

	// Make sure the table isn't modified concurrently
	table.Mutex.Lock()
	defer table.Mutex.Unlock()

	if _, ok := table.Members.Load(client); ok {
		return localization.ErrorTableAlreadyJoined
	}
	table.Members.Store(client, &TableMember{
		Client:  client,
		Color:   color,
		Enabled: false,
	})
	table.MemberCount++

	return nil
}

// Change the enabled state for a member
func ChangeTableMemberState(room string, client string, enabled bool) localization.Translations {

	// Get the table
	obj, valid := tablesCache.Load(room)
	if !valid {
		return localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	// Make sure the table isn't modified concurrently
	table.Mutex.Lock()
	defer table.Mutex.Unlock()

	// Get the member
	obj, valid = table.Members.Load(client)
	if !valid {
		return localization.ErrorTableClientNotFound
	}
	member := obj.(*TableMember)
	member.Enabled = enabled

	return nil
}

func GetTable(room string) (bool, *TableData) {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return false, nil
	}
	return true, obj.(*TableData)
}

func LeaveTable(room string, client string) localization.Translations {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	// Delete the client from the members list
	table.Mutex.Lock()
	table.Members.Delete(client)
	table.MemberCount--
	table.Mutex.Unlock()

	// If it was the last one, close the table
	if table.MemberCount <= 0 {
		tablesCache.Delete(room)
	}

	return nil
}

type TableObject struct {
	Mutex             *sync.Mutex `json:"-"`
	ID                string      `json:"id"`
	Order             uint        `json:"o"` // Drawing order
	LocationX         float64     `json:"x"`
	LocationY         float64     `json:"y"`
	Width             float64     `json:"w"`
	Height            float64     `json:"h"`
	Rotation          float64     `json:"r"`
	Type              int         `json:"t"`
	Creator           string      `json:"cr"` // ID of the creator
	Holder            string      `json:"ho"` // ID of the current card holder (others can't move/modify it while it's held)
	ModificationQueue []string    `json:"-"`  // Queue of holders wanting to interact with the object
	Data              string      `json:"d"`  // Encrypted
}

// * Object helpers
func AddObjectToTable(room string, object *TableObject) localization.Translations {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	table.Mutex.Lock()

	// Generate random and unique id
	id := util.GenerateToken(6)
	for {
		_, valid := table.Objects.Load(id)
		if !valid {
			break
		}
		id = util.GenerateToken(6)
	}

	// Put object into cache and list
	object.ID = id
	table.Objects.Store(id, object)

	// Give the object the highest order
	if table.highestObject == nil {
		object.Order = 1
		table.highestObject = object
	} else {
		object.Order = table.highestObject.Order + 1
		table.highestObject = object
	}

	table.Mutex.Unlock()

	return nil
}

func RemoveObjectFromTable(room string, object string) localization.Translations {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	// Remove object from cache and list
	table.Objects.Delete(object)

	// Make sure there are no concurrent reads/writes on highestObject
	table.Mutex.Lock()
	defer table.Mutex.Unlock()

	// Set the highest object to a new one in case the highest object is being deleted
	if table.highestObject.ID == object {

		// Get the second highest object
		objId := "-"
		maxOrder := uint(0)
		table.Objects.Range(func(key, value any) bool {
			obj := value.(*TableObject)

			// Prevent object from being modified at the same time
			obj.Mutex.Lock()
			defer obj.Mutex.Unlock()

			if obj.Order > maxOrder {
				maxOrder = obj.Order
				objId = obj.ID
			}

			return true
		})

		if objId != "-" {
			// Set the object as the new highest
			return MarkAsNewHighest(room, objId, true, false)
		} else {
			// If it's the last object on the table, set the highest object to nil
			table.highestObject = nil
		}
	}

	return nil
}

func ModifyTableObject(room string, client string, objectId string, data string, width float64, height float64) localization.Translations {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	// Modify object data
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return localization.ErrorObjectNotFound
	}
	object := tObj.(*TableObject)

	// Prevent object from being modified at the same time
	object.Mutex.Lock()
	defer object.Mutex.Unlock()

	// Check if the client is in the modification queue
	if object.ModificationQueue[0] != client {
		return localization.ErrorObjectNotInQueue
	}

	// Modify the data and stuff
	object.Data = data
	object.Width = width
	object.Height = height

	return nil
}

func MoveTableObject(room string, client string, objectId string, x, y float64) localization.Translations {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	// Modify object data
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return localization.ErrorObjectNotFound
	}
	object := tObj.(*TableObject)

	// Check if the client is actually holding the object
	if object.Holder != client {
		return localization.ErrorObjectAlreadyHeld
	}

	// Prevent object from being modified at the same time
	object.Mutex.Lock()
	defer object.Mutex.Unlock()

	object.LocationX = x
	object.LocationY = y

	return nil
}

func RotateTableObject(room string, client string, objectId string, rotation float64) localization.Translations {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	// Modify object data
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return localization.ErrorObjectNotFound
	}
	object := tObj.(*TableObject)

	// Check if the client is in the modification queue
	if object.ModificationQueue[0] != client {
		return localization.ErrorObjectNotInQueue
	}

	// Prevent object from being modified at the same time
	object.Mutex.Lock()
	defer object.Mutex.Unlock()

	// Change the rotation
	object.Rotation = rotation

	return nil
}

func TableObjects(room string) ([]*TableObject, localization.Translations) {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return nil, localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	// Get all the objects from the map
	objects := []*TableObject{}
	table.Objects.Range(func(key, value any) bool {
		objects = append(objects, value.(*TableObject))
		return true
	})

	return objects, nil
}

// Select a table object (no-one else will be able to modify it)
func SelectTableObject(room string, objectId string, client string) localization.Translations {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	// Load the object
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return localization.ErrorObjectNotFound
	}
	object := tObj.(*TableObject)

	// Set the new holder, if possible
	if object.Holder != "" {
		return localization.ErrorObjectAlreadyHeld
	}

	// Prevent object from being modified at the same time
	object.Mutex.Lock()
	defer object.Mutex.Unlock()

	// Set the actual holder
	object.Holder = client

	return nil
}

// Unselect a table object
func UnselectTableObject(room string, objectId string, client string) localization.Translations {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	// Load the object
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return localization.ErrorObjectNotFound
	}
	object := tObj.(*TableObject)

	// Prevent object from being modified at the same time
	object.Mutex.Lock()
	defer object.Mutex.Unlock()

	// Unselect it
	object.Holder = ""

	return nil
}

// Queue a new modification (returns whether the client can modify right away)
func QueueTableObjectModification(room string, objectId string, client string) (bool, localization.Translations) {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return false, localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	// Load the object
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return false, localization.ErrorObjectNotFound
	}
	object := tObj.(*TableObject)

	// Prevent object from being modified at the same time
	object.Mutex.Lock()
	defer object.Mutex.Unlock()

	// Remove all disconnected clients from the queue
	object.ModificationQueue = slices.DeleteFunc(object.ModificationQueue, func(element string) bool {
		_, valid := table.Members.Load(client)
		return !valid
	})

	// Add a new client to the modification queue
	object.ModificationQueue = append(object.ModificationQueue, client)

	// Return whether the client is the only one in the queue
	return len(object.ModificationQueue) == 1, nil
}

// Get the next client queued for modification (at index 1)
func NextModifier(room string, objectId string) (string, localization.Translations) {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return "", localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	// Load the object
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return "", localization.ErrorObjectNotFound
	}
	object := tObj.(*TableObject)

	// Check if there is a new client queued for modification
	if len(object.ModificationQueue) == 0 {
		return "", nil
	}

	// Return the next client
	return object.ModificationQueue[0], nil
}

// Remove the current client modifying from the queue
func RemoveFromModificationQueue(room string, objectId string) localization.Translations {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	// Load the object
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return localization.ErrorObjectNotFound
	}
	object := tObj.(*TableObject)

	// Make sure the object isn't modified at the same time
	object.Mutex.Lock()
	defer object.Mutex.Unlock()

	// Remove the current modifier at index 0
	object.ModificationQueue = object.ModificationQueue[1:len(object.ModificationQueue)]

	return nil
}

func GetTableObject(room string, objectId string) (*TableObject, bool) {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return nil, false
	}
	table := obj.(*TableData)

	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return nil, false
	}
	return tObj.(*TableObject), true
}

func GetMemberData(room string, connId string) (*TableMember, bool) {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return nil, false
	}
	table := obj.(*TableData)

	tObj, valid := table.Members.Load(connId)
	if !valid {
		return nil, false
	}
	return tObj.(*TableMember), true
}

// Mark an object on the table as the new highest (also notifies clients about it).
//
// Set excludeLast to true, if you don't want the client to know about the last object.
// Set lockTable to true, if you want the table mutex to be locked (should be true by default unless
// you know what you're doing)
func MarkAsNewHighest(room string, objectId string, excludeLast bool, lockTable bool) localization.Translations {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	// Load the object
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return localization.ErrorObjectNotFound
	}
	object := tObj.(*TableObject)

	// Prevent object from being modified at the same time
	object.Mutex.Lock()
	defer object.Mutex.Unlock()

	// Prevent table from being modified at the same time
	if lockTable {
		table.Mutex.Lock()
		defer table.Mutex.Unlock()
	}

	// Make sure the same object isn't swapped
	if object.ID == table.highestObject.ID {
		return nil
	}

	// Make sure there is currently a highest object
	if table.highestObject == nil {
		return nil
	}

	// Prevent highest object from being modified at the same time
	currentHighest := table.highestObject
	currentHighest.Mutex.Lock()
	defer currentHighest.Mutex.Unlock()

	// Mark the object as the new highest object
	lastOrder := object.Order
	lastObject := currentHighest.ID
	object.Order = currentHighest.Order
	if !excludeLast {
		currentHighest.Order = lastOrder
	}
	table.highestObject = object

	// Send an event notifying everyone of the swap
	if excludeLast {
		SendEventToMembers(room, pipes.Event{
			Name: "tobj_order",
			Data: map[string]interface{}{
				"o":   table.highestObject.ID,
				"or":  table.highestObject.Order,
				"lo":  lastObject,
				"lor": -1, // To signal to the client that it was removed
			},
		})
	} else {
		SendEventToMembers(room, pipes.Event{
			Name: "tobj_order",
			Data: map[string]interface{}{
				"o":   table.highestObject.ID,
				"or":  table.highestObject.Order,
				"lo":  lastObject,
				"lor": lastOrder,
			},
		})
	}

	return nil
}

// Send an event to all table members
func SendEventToMembers(room string, event pipes.Event) bool {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return false
	}
	data := obj.(*TableData)

	// Get a list of all the adapters of the members
	adapters := []string{}
	data.Members.Range(func(_, value any) bool {

		// Only add the member if they are part of the table
		member := value.(*TableMember)
		if member.Enabled {
			adapters = append(adapters, member.Client)
			util.Log.Println("adding ", member.Client)
		}
		return true
	})

	util.Log.Println("hi hi hi")

	// Send the event through pipes
	if err := SSNode.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.BroadcastChannel(adapters),
		Local:   true,
		Event:   event,
	}); err != nil {
		util.Log.Println("error during event sending to tabletop members:", err)
		return false
	}

	return true
}
