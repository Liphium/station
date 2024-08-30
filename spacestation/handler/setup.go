package handler

import (
	"context"
	"os"

	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
)

// Action: setup
func setup(c *pipeshandler.Context, data string) pipes.Event {

	// Generate new connection
	connection := caching.EmptyConnection(c.Client.ID, c.Client.Session)

	// Insert data
	if !caching.SetMemberData(c.Client.Session, c.Client.ID, connection.ClientID, data) {
		return pipeshandler.ErrorResponse(c, "invalid", nil)
	}

	// Send the update to all members in the room
	if !SendRoomData(c.Client.Session) {
		return pipeshandler.ErrorResponse(c, integration.ErrorServer, nil)
	}

	// Check if livekit room already exists
	rooms, err := caching.RoomClient.ListRooms(context.Background(), &livekit.ListRoomsRequest{
		Names: []string{c.Client.Session},
	})
	if err != nil {
		return pipeshandler.ErrorResponse(c, integration.ErrorServer, err)
	}

	if len(rooms.Rooms) > 0 {

		// Generate livekit token
		token := caching.RoomClient.CreateToken()
		token.AddGrant(&auth.VideoGrant{
			RoomJoin:          true,
			Room:              c.Client.Session,
			CanPublishSources: []string{"camera", "microphone", "screenshare"},
		})
		token.SetIdentity(connection.ClientID)

		jwtToken, err := token.ToJWT()
		if err != nil {
			return pipeshandler.ErrorResponse(c, integration.ErrorServer, err)
		}

		return pipeshandler.NormalResponse(c, map[string]interface{}{
			"success": true,
			"id":      connection.ClientID,
			"key":     connection.KeyBase64(),
			"port":    util.UDPPort,
			"url":     os.Getenv("SS_LK_URL"),
			"token":   jwtToken,
		})
	}

	util.Log.Println("creating new room for", c.Client.Session)

	_, err = caching.RoomClient.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
		Name:            c.Client.Session,
		EmptyTimeout:    120,
		MaxParticipants: 100,
	})
	if err != nil {
		return pipeshandler.ErrorResponse(c, integration.ErrorServer, err)
	}

	// Generate livekit token
	token := caching.RoomClient.CreateToken()
	token.AddGrant(&auth.VideoGrant{
		RoomJoin: true,
		Room:     c.Client.Session,
	})
	token.SetIdentity(connection.ClientID)
	jwtToken, err := token.ToJWT()
	if err != nil {
		return pipeshandler.ErrorResponse(c, integration.ErrorServer, err)
	}

	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"id":      connection.ClientID,
		"key":     connection.KeyBase64(),
		"url":     os.Getenv("SS_LK_URL"),
		"token":   jwtToken,
	})
}
