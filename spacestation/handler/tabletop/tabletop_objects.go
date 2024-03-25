package tabletop_handlers

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
)

// Action: tobj_create
func createObject(message pipeshandler.Context) {

	if message.ValidateForm("x", "y", "w", "h", "r", "type", "data") {
		pipeshandler.ErrorResponse(message, "invalid")
		return
	}

	connection, valid := caching.GetConnection(message.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(message, "invalid")
		return
	}

	x := message.Data["x"].(float64)
	y := message.Data["y"].(float64)
	width := message.Data["w"].(float64)
	height := message.Data["h"].(float64)
	rotation := message.Data["r"].(float64)
	objType := int(message.Data["type"].(float64))
	objData := message.Data["data"].(string)

	object := &caching.TableObject{
		LocationX: x,
		LocationY: y,
		Width:     width,
		Height:    height,
		Rotation:  rotation,
		Type:      objType,
		Data:      objData,
		Creator:   message.Client.ID,
	}
	err := caching.AddObjectToTable(message.Client.Session, object)
	if err != nil {
		pipeshandler.ErrorResponse(message, "server.error")
		return
	}

	// Notify other clients
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
		"id":      object.ID,
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
		pipeshandler.ErrorResponse(ctx, "server.error")
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
		pipeshandler.ErrorResponse(ctx, util.ErrorTabletopInvalidAction)
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

	connection, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Check if the object is held by the client
	obj, valid := caching.GetTableObject(ctx.Client.Session, ctx.Data["id"].(string))
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}
	if obj.Holder != connection.ClientID {
		pipeshandler.ErrorResponse(ctx, util.ErrorTabletopInvalidAction)
		return
	}

	err := caching.ModifyTableObject(ctx.Client.Session, ctx.Data["id"].(string), ctx.Data["data"].(string),
		ctx.Data["width"].(float64), ctx.Data["height"].(float64))
	if err != nil {
		pipeshandler.ErrorResponse(ctx, "server.error")
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

	pipeshandler.SuccessResponse(ctx)
}

// Action: tobj_move
func moveObject(ctx pipeshandler.Context) {

	if ctx.ValidateForm("id", "x", "y") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	x := ctx.Data["x"].(float64)
	y := ctx.Data["y"].(float64)

	err := caching.MoveTableObject(ctx.Client.Session, ctx.Data["id"].(string), x, y)
	if err != nil {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	// Notify other clients
	valid := SendEventToMembers(ctx.Client.Session, pipes.Event{
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

	if ctx.ValidateForm("id", "r") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	connection, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	rotation := ctx.Data["r"].(float64)

	err := caching.RotateTableObject(ctx.Client.Session, ctx.Data["id"].(string), rotation)
	if err != nil {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	// Notify other clients
	valid = SendEventToMembers(ctx.Client.Session, pipes.Event{
		Name: "tobj_rotated",
		Data: map[string]interface{}{
			"id": ctx.Data["id"].(string),
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
