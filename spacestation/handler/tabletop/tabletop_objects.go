package tabletop_handlers

import (
	"sync"

	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
)

// Action: tobj_create
func createObject(message pipeshandler.Context) {

	// Validate message integrity
	if message.ValidateForm("x", "y", "w", "h", "r", "type", "data") {
		pipeshandler.ErrorResponse(message, "invalid")
		return
	}

	// Get the connection (for the client id)
	connection, valid := caching.GetConnection(message.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(message, "invalid")
		return
	}

	// Get all data from the message
	x := message.Data["x"].(float64)
	y := message.Data["y"].(float64)
	width := message.Data["w"].(float64)
	height := message.Data["h"].(float64)
	rotation := message.Data["r"].(float64)
	objType := int(message.Data["type"].(float64))
	objData := message.Data["data"].(string)

	// Create the object here so the data is still there when we send it down below
	object := &caching.TableObject{
		Mutex:             &sync.Mutex{},
		LocationX:         x,
		LocationY:         y,
		Width:             width,
		Height:            height,
		Rotation:          rotation,
		Type:              objType,
		Data:              objData,
		Creator:           message.Client.ID,
		Holder:            "",
		ModificationQueue: []string{},
	}

	// Add the object to the table
	err := caching.AddObjectToTable(message.Client.Session, object)
	if err != nil {
		pipeshandler.ErrorResponse(message, err.Error())
		return
	}

	// Notify other clients about the object creation
	valid = SendEventToMembers(message.Client.Session, pipes.Event{
		Name: "tobj_created",
		Data: map[string]interface{}{
			"id":   object.ID,
			"x":    x,
			"y":    y,
			"w":    width,
			"h":    height,
			"r":    rotation,
			"type": objType,
			"data": objData,
			"c":    connection.ClientID,
		},
	})
	if !valid {
		pipeshandler.ErrorResponse(message, "server.error")
		return
	}

	pipeshandler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"id":      object.ID, // So the client can set the new id
	})
}

// Action: tobj_delete
func deleteObject(ctx pipeshandler.Context) {

	if ctx.ValidateForm("id") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	err := caching.RemoveObjectFromTable(ctx.Client.Session, ctx.Data["id"].(string))
	if err != nil {
		pipeshandler.ErrorResponse(ctx, err.Error())
		return
	}

	// Notify other clients
	valid := SendEventToMembers(ctx.Client.Session, pipes.Event{
		Name: "tobj_deleted",
		Data: map[string]interface{}{
			"id": ctx.Data["id"].(string),
		},
	})
	if !valid {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	pipeshandler.SuccessResponse(ctx)
}

// Action: tobj_select
func selectObject(ctx pipeshandler.Context) {

	if ctx.ValidateForm("id") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	connection, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Grab hold of it
	err := caching.SelectTableObject(ctx.Client.Session, ctx.Data["id"].(string), connection.ClientID)
	if err != nil {
		pipeshandler.ErrorResponse(ctx, err.Error())
		return
	}

	pipeshandler.SuccessResponse(ctx)
}

// Action: tobj_select
func unselectObject(ctx pipeshandler.Context) {

	// Validate message integrity
	if ctx.ValidateForm("id") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Get the connection (for the client id)
	connection, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Grab hold of it
	err := caching.UnselectTableObject(ctx.Client.Session, ctx.Data["id"].(string), connection.ClientID)
	if err != nil {
		pipeshandler.ErrorResponse(ctx, err.Error())
		return
	}

	pipeshandler.SuccessResponse(ctx)
}

// Action: tobj_modify
func modifyObject(ctx pipeshandler.Context) {

	if ctx.ValidateForm("id", "data", "width", "height") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Get the connection (for the client id)
	connection, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Get the object id from the message
	objectId := ctx.Data["id"].(string)

	// Make sure the next client gets to modify the object regardless of errors
	defer handleNextModification(ctx.Client.Session, objectId)

	// Modify the object and return the error if there is one
	err := caching.ModifyTableObject(ctx.Client.Session, connection.ClientID, objectId, ctx.Data["data"].(string),
		ctx.Data["width"].(float64), ctx.Data["height"].(float64))
	if err != nil {
		pipeshandler.ErrorResponse(ctx, err.Error())
		return
	}

	// Notify other clients
	util.Log.Println("Sending tobj_modified event")
	valid = SendEventToMembers(ctx.Client.Session, pipes.Event{
		Name: "tobj_modified",
		Data: map[string]interface{}{
			"id":   ctx.Data["id"].(string),
			"data": ctx.Data["data"].(string),
			"w":    ctx.Data["width"].(float64),
			"h":    ctx.Data["height"].(float64),
		},
	})
	if !valid {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	// Add the next client to the modification queue

	pipeshandler.SuccessResponse(ctx)
}

// Action: tobj_move
func moveObject(ctx pipeshandler.Context) {

	// Validate message integrity
	if ctx.ValidateForm("id", "x", "y") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Get the connection (for the client id)
	connection, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Get data from the message
	x := ctx.Data["x"].(float64)
	y := ctx.Data["y"].(float64)

	// Move the actual object
	err := caching.MoveTableObject(ctx.Client.Session, connection.ClientID, ctx.Data["id"].(string), x, y)
	if err != nil {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	// Notify other clients
	valid = SendEventToMembers(ctx.Client.Session, pipes.Event{
		Name: "tobj_moved",
		Data: map[string]interface{}{
			"id": ctx.Data["id"].(string),
			"x":  x,
			"y":  y,
		},
	})
	if !valid {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	pipeshandler.SuccessResponse(ctx)
}

// Action: tobj_rotate
func rotateObject(ctx pipeshandler.Context) {

	// Validate message integrity
	if ctx.ValidateForm("id", "r") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Get the connection (for the client id)
	connection, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Get the data from the message
	objectId := ctx.Data["id"].(string)
	rotation := ctx.Data["r"].(float64)

	// Make sure the next client gets to modify the object regardless of errors
	defer handleNextModification(ctx.Client.Session, objectId)

	// Rotate the object and return an error (only if one is there)
	err := caching.RotateTableObject(ctx.Client.Session, connection.ClientID, objectId, rotation)
	if err != nil {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	// Notify other clients about the rotation
	valid = SendEventToMembers(ctx.Client.Session, pipes.Event{
		Name: "tobj_rotated",
		Data: map[string]interface{}{
			"id": objectId,
			"s":  connection.ClientID,
			"r":  rotation,
		},
	})
	if !valid {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	pipeshandler.SuccessResponse(ctx)
}

// Action: tobj_mqueue
func queueModificationToObject(ctx pipeshandler.Context) {

	// Validate message integrity
	if ctx.ValidateForm("id") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Get the connection (for the client id)
	connection, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Get the id of the object from the message
	objectId := ctx.Data["id"].(string)

	// Queue the modification
	rightAway, err := caching.QueueTableObjectModification(ctx.Client.Session, objectId, connection.ClientID)
	if err != nil {
		pipeshandler.ErrorResponse(ctx, err.Error())
		return
	}

	// Return whether the modification can be sent right away
	pipeshandler.NormalResponse(ctx, map[string]interface{}{
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
