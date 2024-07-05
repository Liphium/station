package server

import (
	"context"
	"log"
	"os"

	"github.com/livekit/protocol/auth"
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

	// Create some test tokens
	for _, name := range []string{"test1", "test2"} {
		token := RoomClient.CreateToken()
		token.AddGrant(&auth.VideoGrant{
			RoomJoin: true,
			Room:     "test",
		})
		token.SetIdentity(name)

		jwtToken, err := token.ToJWT()
		if err != nil {
			log.Println("shit")
			return
		}

		log.Println(name + ":" + jwtToken)
	}

}
