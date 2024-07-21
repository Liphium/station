package caching

import (
	"errors"
	"slices"
	"sync"

	"github.com/Liphium/station/spacestation/util"
)

// ! For setting please ALWAYS use cost 1
// Room ID -> Table
var tablesCache *sync.Map = &sync.Map{}

// Errors
var (
	ErrTableNotFound            = errors.New("tabletop.not_found")
	ErrClientAlreadyJoinedTable = errors.New("tabletop.already_joined")
	ErrCouldntCreateTable       = errors.New("tabletop.couldnt_create")
	ErrObjectNotFound           = errors.New("tabletop.object_not_found")
	ErrObjectAlreadyHeld        = errors.New("tabletop.object_already_held")
	ErrObjectDifferentHolder    = errors.New("tabletop.object_different_holder")
)

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
}

// * Table management
func JoinTable(room string, client string, color float64) error {

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

	table.Mutex.Lock()
	if _, ok := table.Members.Load(client); ok {
		return ErrClientAlreadyJoinedTable
	}
	table.Members.Store(client, &TableMember{
		Client: client,
		Color:  color,
	})
	table.MemberCount++
	table.Mutex.Unlock()

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

func LeaveTable(room string, client string) error {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return ErrTableNotFound
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
	ID        string  `json:"id"`
	LocationX float64 `json:"x"`
	LocationY float64 `json:"y"`
	Width     float64 `json:"w"`
	Height    float64 `json:"h"`
	Rotation  float64 `json:"r"`
	Type      int     `json:"t"`
	Creator   string  `json:"cr"` // ID of the creator
	Holder    string  `json:"ho"` // ID of the current card holder (others can't move/modify it while it's held)
	Data      string  `json:"d"`  // Encrypted
}

// * Object helpers
func AddObjectToTable(room string, object *TableObject) error {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return ErrTableNotFound
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

func RemoveObjectFromTable(room string, object string) error {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return ErrTableNotFound
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

func ModifyTableObject(room string, client string, objectId string, data string, width float64, height float64) error {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return ErrTableNotFound
	}
	table := obj.(*TableData)

	// Modify object data
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return ErrObjectNotFound
	}
	object := tObj.(*TableObject)

	// Check if the holder is the client
	if object.Holder != client {
		return ErrObjectDifferentHolder
	}

	object.Data = data
	object.Width = width
	object.Height = height

	return nil
}

func MoveTableObject(room string, client string, objectId string, x, y float64) error {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return ErrTableNotFound
	}
	table := obj.(*TableData)

	// Modify object data
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return ErrObjectNotFound
	}
	object := tObj.(*TableObject)

	// Check if the client is actually holding the object
	if object.Holder != client {
		return ErrObjectDifferentHolder
	}

	object.LocationX = x
	object.LocationY = y

	return nil
}

func RotateTableObject(room string, client string, objectId string, rotation float64) error {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return ErrTableNotFound
	}
	table := obj.(*TableData)

	// Modify object data
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return ErrObjectNotFound
	}
	object := tObj.(*TableObject)

	// Check if the client is actually holding the object
	if object.Holder != client {
		return ErrObjectDifferentHolder
	}

	// Change the rotation
	object.Rotation = rotation

	return nil
}

func TableObjects(room string) ([]*TableObject, error) {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return nil, ErrTableNotFound
	}
	table := obj.(*TableData)

	objects := make([]*TableObject, len(table.ObjectList))
	for i, value := range table.ObjectList {
		object, valid := table.Objects.Load(value)
		if !valid {
			return nil, ErrObjectNotFound
		}

		objects[i] = object.(*TableObject)
	}

	return objects, nil
}

// Select a table object (no-one else will be able to modify it)
func SelectTableObject(room string, objectId string, client string) error {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return ErrTableNotFound
	}
	table := obj.(*TableData)

	// Load the object
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return ErrObjectNotFound
	}
	object := tObj.(*TableObject)

	// Set the new holder, if possible
	if object.Holder != "" {
		return ErrObjectAlreadyHeld
	}
	object.Holder = client

	return nil
}

// Unselect a table object
func UnselectTableObject(room string, objectId string, client string) error {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return ErrTableNotFound
	}
	table := obj.(*TableData)

	// Load the object
	tObj, valid := table.Objects.Load(objectId)
	if !valid {
		return ErrObjectNotFound
	}
	object := tObj.(*TableObject)

	// Unselect it
	object.Holder = ""

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

func GetMemberData(room string, client string) (*TableMember, bool) {
	obj, valid := tablesCache.Load(room)
	if !valid {
		return nil, false
	}
	table := obj.(*TableData)

	tObj, valid := table.Members.Load(client)
	if !valid {
		return nil, false
	}
	return tObj.(*TableMember), true
}
