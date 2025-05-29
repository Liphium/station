package util

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

// Environment variables
const EnvAppName = "APP_NAME" // Configure the app name

// Locals constants
const LocalsServerPriv = "srv_priv"
const LocalsServerPub = "srv_pub"
const LocalsKey = "key"
const LocalsBody = "body"

// Important variables
const ProtocolVersion = 8

var Testing = false
var LogErrors = true

var Log = log.New(os.Stdout, "backend ", log.Flags())

var JWT_SECRET = ""

var NodeProtocol = "http://"

// Send a post request
func PostRequest(url string, body map[string]interface{}) (map[string]interface{}, error) {

	req, _ := sonic.Marshal(body)

	reader := strings.NewReader(string(req))

	res, err := http.Post(url, "application/json", reader)
	if err != nil {
		return nil, err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, res.Body)

	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = sonic.Unmarshal([]byte(buf.String()), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Get the system uuid set in the environment variables
func GetSystemUUID() uuid.UUID {
	id, err := uuid.Parse(os.Getenv("SYSTEM_UUID"))
	if err != nil {
		panic("Please set the SYSTEM_UUID env property.")
	}

	return id
}
