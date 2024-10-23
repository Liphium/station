package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/Liphium/station/main/integration"
)

// The domain of the chat server
var OwnPath = ""

const StatusOnline = 0
const StatusOffline = 1
const StatusError = 2

var Log = log.New(os.Stdout, "chat-server ", log.Flags())
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateToken(tkLength int32) string {

	s := make([]rune, tkLength)

	length := big.NewInt(int64(len(letters)))

	for i := range s {

		number, _ := rand.Int(rand.Reader, length)
		s[i] = letters[number.Int64()]
	}

	return string(s)
}

// Hashes using SHA256
func HashString(str string) string {

	hashed := sha256.Sum256([]byte(str))
	return base64.StdEncoding.EncodeToString(hashed[:])
}

// Hashes using SHA256
func Hash(str string) []byte {

	hashed := sha256.Sum256([]byte(str))
	return hashed[:]
}

func Node64(id int64) string {
	return fmt.Sprintf("%d", id)
}

func NodeTo64(id string) int64 {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0
	}

	return i
}

func PrintIfTesting(v ...any) {
	if integration.Testing {
		Log.Println(v)
	}
}
