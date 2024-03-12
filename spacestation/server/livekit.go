package server

import (
	"context"
	"os"

	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go"
)

var RoomClient *lksdk.RoomServiceClient

func InitLiveKit() {
	RoomClient = lksdk.NewRoomServiceClient(os.Getenv("SS_LK_URL"), os.Getenv("SS_LK_KEY"), os.Getenv("SS_LK_SECRET"))

	_, err := RoomClient.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
		Name:         "test",
		EmptyTimeout: 60,
	})
	if err != nil {
		panic(err)
	}
}
