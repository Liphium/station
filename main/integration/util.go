package integration

import (
	"bytes"
	"crypto/rand"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
)

var Testing = false
var TestingToken = ""

var JwtSecret = ""

const StatusOnline = 0
const StatusOffline = 1
const StatusError = 2

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

var Protocol = "http://"
var BasePath = "http://localhost:3000"
var Domain = "localhost:3000"

// * Important
const ApiVersion = "v1"

// Send a post request (without TC encryption and custom URL)
func PostRequestURL(url string, body map[string]interface{}) (map[string]interface{}, error) {

	byteBody, err := sonic.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader := strings.NewReader(string(byteBody))

	res, err := http.Post(url, "application/json", bodyReader)
	if err != nil {
		return nil, err
	}

	// Decrypt the request body
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return nil, err
	}

	// Parse decrypted body into JSON
	var data map[string]interface{}
	err = sonic.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Send a post request
func PostRequestBackend(url string, body map[string]interface{}) (map[string]interface{}, error) {
	return PostRequest(BasePath, "/"+ApiVersion+url, body)
}

// Send a post request
func PostRequestBackendGeneric[T any](url string, body map[string]interface{}) (T, error) {
	return PostRequestGeneric[T](BasePath, "/"+ApiVersion+url, body)
}

// Send a post request
func PostRequestBackendServer(server string, url string, body map[string]interface{}) (map[string]interface{}, error) {
	return PostRequest(server, "/"+ApiVersion+url, body)
}

// Send a post request
func PostRequestBackendServerGeneric[T any](server string, url string, body map[string]interface{}) (T, error) {
	return PostRequestGeneric[T](server, "/"+ApiVersion+url, body)
}

// Send a post request (no generics)
func PostRequest(server string, path string, body map[string]interface{}) (map[string]interface{}, error) {
	return PostRequestGeneric[map[string]interface{}](server, path, body)
}

// Send a post request
func PostRequestGeneric[T any](server string, path string, body map[string]interface{}) (T, error) {

	// Declared here so it can be returned as nil before it's actually used
	var data T

	// Make sure there is a protocol specified on the server
	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		server = "https://" + server
	}

	// Encode body to JSON
	byteBody, err := sonic.Marshal(body)
	if err != nil {
		return data, err
	}

	// Set headers
	reqHeaders := http.Header{}
	reqHeaders.Set("Content-Type", "application/json")

	// Send the request
	req, err := http.NewRequest(http.MethodPost, server+path, bytes.NewBuffer(byteBody))
	if err != nil {
		return data, err
	}
	req.Header = reqHeaders

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return data, err
	}

	// Grab all bytes from the buffer
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return data, err
	}

	// Parse body into JSON
	err = sonic.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

type NormalResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
