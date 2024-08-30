package tabletop_handlers

import (
	"sync"

	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
)

// Action: tobj_create
func createObject(c *pipeshandler.Context, action struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Width    float64 `json:"w"`
	Height   float64 `json:"h"`
	Rotation float64 `json:"r"`
	Type     int     `json:"type"`
	Data     string  `json:"data"`
}) pipes.Event {

	// Get the connection (for the client id)
	connection, valid := caching.GetConnection(c.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(c, "invalid", nil)
	}

	// Create the object here so the data is still there when we send it down below
	object := &caching.TableObject{
		Mutex:             &sync.Mutex{},
		LocationX:         action.X,
		LocationY:         action.Y,
		Width:             action.Width,
		Height:            action.Height,
		Rotation:          action.Rotation,
		Type:              action.Type,
		Data:              action.Data,
		Creator:           c.Client.ID,
		Holder:            "",
		ModificationQueue: []string{},
	}

	// Add the object to the table
	err := caching.AddObjectToTable(c.Client.Session, object)
	if err != nil {
		return pipeshandler.ErrorResponse(c, err.Error(), err)
	}

	// Notify other clients about the object creation
	valid = SendEventToMembers(c.Client.Session, pipes.Event{
		Name: "tobj_created",
		Data: map[string]interface{}{
			"id":   object.ID,
			"x":    action.X,
			"y":    action.Y,
			"w":    action.Width,
			"h":    action.Height,
			"r":    action.Rotation,
			"type": action.Type,
			"data": action.Data,
			"c":    connection.ClientID,
		},
	})
	if !valid {
		return pipeshandler.ErrorResponse(c, "server.error", nil)
	}

	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"id":      object.ID, // So the client can set the new id
	})
}

// Action: tobj_delete
func deleteObject(c *pipeshandler.Context, id string) pipes.Event {

	err := caching.RemoveObjectFromTable(c.Client.Session, id)
	if err != nil {
		return pipeshandler.ErrorResponse(c, err.Error(), err)
	}

	// Notify other clients
	valid := SendEventToMembers(c.Client.Session, pipes.Event{
		Name: "tobj_deleted",
		Data: map[string]interface{}{
			"id": id,
		},
	})
	if !valid {
		return pipeshandler.ErrorResponse(c, "server.error", nil)
	}

	return pipeshandler.SuccessResponse(c)
}

// Action: tobj_select
func selectObject(c *pipeshandler.Context, id string) pipes.Event {

	connection, valid := caching.GetConnection(c.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(c, "invalid", nil)
	}

	// Grab hold of it
	err := caching.SelectTableObject(c.Client.Session, id, connection.ClientID)
	if err != nil {
		return pipeshandler.ErrorResponse(c, err.Error(), err)
	}

	return pipeshandler.SuccessResponse(c)
}

// Action: tobj_select
func unselectObject(c *pipeshandler.Context, id string) pipes.Event {

	// Get the connection (for the client id)
	connection, valid := caching.GetConnection(c.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(c, "invalid", nil)
	}

	// Grab hold of it
	err := caching.UnselectTableObject(c.Client.Session, id, connection.ClientID)
	if err != nil {
		return pipeshandler.ErrorResponse(c, err.Error(), err)
	}

	return pipeshandler.SuccessResponse(c)
}

// Action: tobj_modify
func modifyObject(c *pipeshandler.Context, action struct {
	ID     string  `json:"id"`
	Data   string  `json:"data"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}) pipes.Event {

	// Get the connection (for the client id)
	connection, valid := caching.GetConnection(c.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(c, "invalid", nil)
	}

	// Make sure the next client gets to modify the object regardless of errors
	defer handleNextModification(c.Client.Session, action.ID)

	// Modify the object and return the error if there is one
	err := caching.ModifyTableObject(c.Client.Session, connection.ClientID, action.ID, action.Data, action.Width, action.Height)
	if err != nil {
		return pipeshandler.ErrorResponse(c, err.Error(), err)
	}

	// Notify other clients
	util.Log.Println("Sending tobj_modified event")
	valid = SendEventToMembers(c.Client.Session, pipes.Event{
		Name: "tobj_modified",
		Data: map[string]interface{}{
			"id":   action.ID,
			"data": action.Data,
			"w":    action.Width,
			"h":    action.Height,
		},
	})
	if !valid {
		return pipeshandler.ErrorResponse(c, "server.error", nil)
	}

	// Add the next client to the modification queue (in the defer above)
	return pipeshandler.SuccessResponse(c)
}

// Action: tobj_move
func moveObject(c *pipeshandler.Context, action struct {
	ID string  `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}) pipes.Event {

	// Get the connection (for the client id)
	connection, valid := caching.GetConnection(c.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(c, "invalid", nil)
	}

	// Move the actual object
	err := caching.MoveTableObject(c.Client.Session, connection.ClientID, action.ID, action.X, action.Y)
	if err != nil {
		return pipeshandler.ErrorResponse(c, err.Error(), err)
	}

	// Notify other clients
	valid = SendEventToMembers(c.Client.Session, pipes.Event{
		Name: "tobj_moved",
		Data: map[string]interface{}{
			"id": action.ID,
			"x":  action.X,
			"y":  action.Y,
		},
	})
	if !valid {
		return pipeshandler.ErrorResponse(c, "server.error", nil)
	}

	return pipeshandler.SuccessResponse(c)
}

// Action: tobj_rotate
func rotateObject(c *pipeshandler.Context, action struct {
	ID       string  `json:"id"`
	Rotation float64 `json:"r"`
}) pipes.Event {

	// Get the connection (for the client id)
	connection, valid := caching.GetConnection(c.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(c, "invalid", nil)
	}

	// Make sure the next client gets to modify the object regardless of errors
	defer handleNextModification(c.Client.Session, action.ID)

	// Rotate the object and return an error (only if one is there)
	err := caching.RotateTableObject(c.Client.Session, connection.ClientID, action.ID, action.Rotation)
	if err != nil {
		return pipeshandler.ErrorResponse(c, err.Error(), err)
	}

	// Notify other clients about the rotation
	valid = SendEventToMembers(c.Client.Session, pipes.Event{
		Name: "tobj_rotated",
		Data: map[string]interface{}{
			"id": action.ID,
			"s":  connection.ClientID,
			"r":  action.Rotation,
		},
	})
	if !valid {
		return pipeshandler.ErrorResponse(c, "server.error", nil)
	}

	return pipeshandler.SuccessResponse(c)
}

// Action: tobj_mqueue
func queueModificationToObject(c *pipeshandler.Context, objectId string) pipes.Event {

	// Get the connection (for the client id)
	connection, valid := caching.GetConnection(c.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(c, "invalid", nil)
	}

	// Queue the modification
	rightAway, err := caching.QueueTableObjectModification(c.Client.Session, objectId, connection.ClientID)
	if err != nil {
		return pipeshandler.ErrorResponse(c, err.Error(), err)
	}

	// Return whether the modification can be sent right away
	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"direct":  rightAway,
	})
}

// Called when a modification is completed to contact the next modifier
func handleNextModification(room string, object string) {
	// Remove the current client from the modification queue
	err := caching.RemoveFromModificationQueue(room, object)
	if err != nil {
		util.Log.Println("couldn't remove from modification queue:", err)
		return
	}

	// Get the next client to modify the object
	client, err := caching.NextModifier(room, object)
	if err != nil {
		util.Log.Println("error during getting next modifier:", err)
		return
	}

	// Return if there is no client
	if client == "" {
		return
	}

	// Send event to inform the client about their modification being allowed
	err = caching.SSNode.SendClient(client, pipes.Event{
		Name: "tobj_mqueue_allowed",
		Data: map[string]interface{}{
			"id": object,
		},
	})
	if err != nil {
		util.Log.Println("error during sending next modification:", err)
		return
	}
}
