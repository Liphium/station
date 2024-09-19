package caching

import (
	"slices"
	"sync"

	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/spacestation/util"
)

// ! For setting please ALWAYS use cost 1
// Room ID -> Table
var tablesCache *sync.Map = &sync.Map{}

type TableData struct {
	Mutex       *sync.Mutex
	Room        string
	MemberCount int
	Members     *sync.Map // Client ID -> Client info
	ObjectList  []string  // List of all object ids
	Objects     *sync.Map // Cache for all objects on the table (Object ID -> Object)
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

// rangeFunc returns whether or not the loop should be continued
func RangeOverTableMembers(room string, rangeFunc func(*TableMember) bool) bool {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return false
	}
	data := obj.(*TableData)

	// Range over all members
	data.Members.Range(func(_, value any) bool {
		member := value.(*TableMember)
		return rangeFunc(member)
	})

	return true
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
	id := util.GenerateToken(5)
	for slices.Contains(table.ObjectList, id) {
		id = util.GenerateToken(5)
	}

	// Put object into cache and list
	table.ObjectList = append(table.ObjectList, id)
	object.ID = id
	table.Objects.Store(id, object)

	table.Mutex.Unlock()

	return nil
}

func RemoveObjectFromTable(room string, object string) localization.Translations {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return localization.ErrorTableNotFound
	}
	table := obj.(*TableData)

	table.Mutex.Lock()

	// Put object into cache and list
	for i, member := range table.ObjectList {
		if member == object {
			table.ObjectList = append(table.ObjectList[:i], table.ObjectList[i+1:]...)
			break
		}
	}
	table.Objects.Delete(object)

	table.Mutex.Unlock()

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

	objects := make([]*TableObject, len(table.ObjectList))
	for i, value := range table.ObjectList {
		object, valid := table.Objects.Load(value)
		if !valid {
			return nil, localization.ErrorObjectNotFound
		}

		objects[i] = object.(*TableObject)
	}

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
