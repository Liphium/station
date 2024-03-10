package calls

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	lksdk "github.com/livekit/server-sdk-go"
)

var RoomClient *lksdk.RoomServiceClient

func Connect() {
	RoomClient = lksdk.NewRoomServiceClient(os.Getenv("LK_URL"), os.Getenv("LK_KEY"), os.Getenv("LK_SECRET"))
}

type CallClaims struct {
	CID string `json:"cid"` // Call ID
	Ow  string `json:"ow"`  // Owner ID
	EXP int64  `json:"e_u"` // Expiration time (Unix)
	jwt.RegisteredClaims
}

func GenerateCallToken(id string, owner string) (string, error) {

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, CallClaims{
		CID: id,
		Ow:  owner,
		EXP: time.Now().Add(5 * time.Minute).Unix(),
	})

	token, err := tk.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return token, nil
}

func GetCallClaims(certificate string) (*CallClaims, bool) {

	token, err := jwt.ParseWithClaims(certificate, &CallClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithLeeway(5*time.Minute))

	if err != nil {
		return &CallClaims{}, false
	}

	if claims, ok := token.Claims.(*CallClaims); ok && token.Valid {
		return claims, true
	}

	return &CallClaims{}, false
}

func (m *CallClaims) Valid(id string) bool {
	return m.CID == id && m.EXP > time.Now().Unix()
}
