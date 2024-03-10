package calls

import (
	"os"
	"time"

	"github.com/livekit/protocol/auth"
)

func GetJoinToken(room, identity string) (string, error) {

	// Generate token
	at := auth.NewAccessToken(os.Getenv("LK_KEY"), os.Getenv("LK_SECRET"))

	// Add permissions
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     room,
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	return at.ToJWT()
}
