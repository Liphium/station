package handler

import (
	"context"
	"os"

	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
)

// Action: setup
func setup(ctx pipeshandler.Context) {

	if ctx.ValidateForm("data") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}
	data := ctx.Data["data"].(string)

	// Generate new connection
	connection := caching.EmptyConnection(ctx.Client.ID, ctx.Client.Session)

	// Insert data
	if !caching.SetMemberData(ctx.Client.Session, ctx.Client.ID, connection.ClientID, data) {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	if !SendRoomData(ctx.Client.Session) {
		pipeshandler.ErrorResponse(ctx, integration.ErrorServer)
		return
	}

	// Check if livekit room already exists
	rooms, err := caching.RoomClient.ListRooms(context.Background(), &livekit.ListRoomsRequest{
		Names: []string{ctx.Client.Session},
	})
	if err != nil {
		pipeshandler.ErrorResponse(ctx, integration.ErrorServer)
		return
	}

	if len(rooms.Rooms) > 0 {

		// Generate livekit token
		token := caching.RoomClient.CreateToken()
		token.AddGrant(&auth.VideoGrant{
			RoomJoin:          true,
			Room:              ctx.Client.Session,
			CanPublishSources: []string{"camera", "microphone", "screenshare"},
		})
		token.SetIdentity(connection.ClientID)

		jwtToken, err := token.ToJWT()
		if err != nil {
			pipeshandler.ErrorResponse(ctx, integration.ErrorServer)
			return
		}

		pipeshandler.NormalResponse(ctx, map[string]interface{}{
			"success": true,
			"id":      connection.ClientID,
			"key":     connection.KeyBase64(),
			"port":    util.UDPPort,
			"url":     os.Getenv("SS_LK_URL"),
			"token":   jwtToken,
		})
		return
	}

	util.Log.Println("creating new room for", ctx.Client.Session)

	_, err = caching.RoomClient.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
		Name:            ctx.Client.Session,
		EmptyTimeout:    120,
		MaxParticipants: 100,
	})
	if err != nil {
		pipeshandler.ErrorResponse(ctx, integration.ErrorServer)
		return
	}

	// Generate livekit token
	token := caching.RoomClient.CreateToken()
	token.AddGrant(&auth.VideoGrant{
		RoomJoin: true,
		Room:     ctx.Client.Session,
	})
	token.SetIdentity(connection.ClientID)
	jwtToken, err := token.ToJWT()
	if err != nil {
		pipeshandler.ErrorResponse(ctx, integration.ErrorServer)
		return
	}

	pipeshandler.NormalResponse(ctx, map[string]interface{}{
		"success": true,
		"id":      connection.ClientID,
		"key":     connection.KeyBase64(),
		"url":     os.Getenv("SS_LK_URL"),
		"token":   jwtToken,
	})
}
