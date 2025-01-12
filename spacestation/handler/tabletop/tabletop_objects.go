package tabletop_handlers

import (
	"sync"

	"github.com/Liphium/station/main/localization"
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
	msg := caching.AddObjectToTable(c.Client.Session, object)
	if msg != nil {
		return pipeshandler.ErrorResponse(c, msg, nil)
	}

	// Notify other clients about the object creation
	valid := caching.SendEventToMembers(c.Client.Session, pipes.Event{
		Name: "tobj_created",
		Data: map[string]interface{}{
			"id":   object.ID,
			"x":    action.X,
			"y":    action.Y,
			"o":    object.Order,
			"w":    action.Width,
			"h":    action.Height,
			"r":    action.Rotation,
			"type": action.Type,
			"data": action.Data,
			"c":    c.Client.ID,
		},
	})
	if !valid {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"id":      object.ID,    // So the client can set the new id
		"o":       object.Order, // Cause the client doesn't know order at creation
	})
}

// Action: tobj_delete
func deleteObject(c *pipeshandler.Context, id string) pipes.Event {

	msg := caching.RemoveObjectFromTable(c.Client.Session, id)
	if msg != nil {
		return pipeshandler.ErrorResponse(c, msg, nil)
	}

	// Notify other clients
	valid := caching.SendEventToMembers(c.Client.Session, pipes.Event{
		Name: "tobj_deleted",
		Data: map[string]interface{}{
			"id": id,
		},
	})
	if !valid {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	return pipeshandler.SuccessResponse(c)
}

// Action: tobj_select
func selectObject(c *pipeshandler.Context, id string) pipes.Event {

	// Grab hold of it
	msg := caching.SelectTableObject(c.Client.Session, id, c.Client.ID)
	if msg != nil {
		return pipeshandler.ErrorResponse(c, msg, nil)
	}

	// Move to the highest order
	msg = caching.MarkAsNewHighest(c.Client.Session, id, true)
	if msg != nil {
		return pipeshandler.ErrorResponse(c, msg, nil)
	}

	return pipeshandler.SuccessResponse(c)
}

// Action: tobj_select
func unselectObject(c *pipeshandler.Context, id string) pipes.Event {

	// Grab hold of it
	msg := caching.UnselectTableObject(c.Client.Session, id, c.Client.ID)
	if msg != nil {
		return pipeshandler.ErrorResponse(c, msg, nil)
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

	// Make sure the next client gets to modify the object regardless of errors
	defer handleNextModification(c.Client.Session, action.ID)

	// Modify the object and return the error if there is one
	msg := caching.ModifyTableObject(c.Client.Session, c.Client.ID, action.ID, action.Data, action.Width, action.Height)
	if msg != nil {
		return pipeshandler.ErrorResponse(c, msg, nil)
	}

	// Notify other clients
	util.Log.Println("Sending tobj_modified event")
	if !caching.SendEventToMembers(c.Client.Session, pipes.Event{
		Name: "tobj_modified",
		Data: map[string]interface{}{
			"id":   action.ID,
			"data": action.Data,
			"w":    action.Width,
			"h":    action.Height,
		},
	}) {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
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

	// Move the actual object
	msg := caching.MoveTableObject(c.Client.Session, c.Client.ID, action.ID, action.X, action.Y)
	if msg != nil {
		return pipeshandler.ErrorResponse(c, msg, nil)
	}

	// Notify other clients
	if !caching.SendEventToMembers(c.Client.Session, pipes.Event{
		Name: "tobj_moved",
		Data: map[string]interface{}{
			"id": action.ID,
			"x":  action.X,
			"y":  action.Y,
		},
	}) {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	return pipeshandler.SuccessResponse(c)
}

// Action: tobj_rotate
func rotateObject(c *pipeshandler.Context, action struct {
	ID       string  `json:"id"`
	Rotation float64 `json:"r"`
}) pipes.Event {

	// Make sure the next client gets to modify the object regardless of errors
	defer handleNextModification(c.Client.Session, action.ID)

	// Rotate the object and return an error (only if one is there)
	msg := caching.RotateTableObject(c.Client.Session, c.Client.ID, action.ID, action.Rotation)
	if msg != nil {
		return pipeshandler.ErrorResponse(c, msg, nil)
	}

	// Notify other clients about the rotation
	if !caching.SendEventToMembers(c.Client.Session, pipes.Event{
		Name: "tobj_rotated",
		Data: map[string]interface{}{
			"id": action.ID,
			"s":  c.Client.ID,
			"r":  action.Rotation,
		},
	}) {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	return pipeshandler.SuccessResponse(c)
}

// Action: tobj_mqueue
func queueModificationToObject(c *pipeshandler.Context, objectId string) pipes.Event {

	// Queue the modification
	rightAway, msg := caching.QueueTableObjectModification(c.Client.Session, objectId, c.Client.ID)
	if msg != nil {
		return pipeshandler.ErrorResponse(c, msg, nil)
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
	msg := caching.RemoveFromModificationQueue(room, object)
	if msg != nil {
		util.Log.Println("couldn't remove from modification queue:", msg[localization.DefaultLocale])
		return
	}

	// Get the next client to modify the object
	client, msg := caching.NextModifier(room, object)
	if msg != nil {
		util.Log.Println("error during getting next modifier:", msg[localization.DefaultLocale])
		return
	}

	// Return if there is no client
	if client == "" {
		return
	}

	// Send event to inform the client about their modification being allowed
	err := caching.SSNode.SendClient(client, pipes.Event{
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
