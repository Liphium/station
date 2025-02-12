package util

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Liphium/station/main/localization"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
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
const ProtocolVersion = 7

var Testing = false
var LogErrors = true

var Log = log.New(os.Stdout, "backend ", log.Flags())

var JWT_SECRET = ""

var NodeProtocol = "http://"

// Send a post request (with TC protection encryption)
func PostRequest(key *rsa.PublicKey, url string, body map[string]interface{}) (map[string]interface{}, error) {

	// Encode the json to a byte slice
	byteBody, err := sonic.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Compute the auth tag
	aesKey, err := NewAESKey()
	if err != nil {
		return nil, err
	}
	authTag, err := EncryptRSA(key, aesKey)
	if err != nil {
		return nil, err
	}
	authTagEncoded := base64.StdEncoding.EncodeToString(authTag)

	// Set headers
	reqHeaders := http.Header{}
	reqHeaders.Set("Content-Type", "application/json")
	reqHeaders.Set("Auth-Tag", authTagEncoded)

	// Encrypt the body using the AES key
	encryptedBody, err := EncryptAES(aesKey, byteBody)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(encryptedBody)

	// Create the request and set headers + body
	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, err
	}
	req.Header = reqHeaders

	// Do the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("something went wrong on the node: %d", res.StatusCode)
	}

	// Get the request body in byte slice form
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return nil, err
	}

	// Decrypt the request body byte slice using AES
	decryptedBody, err := DecryptAES(aesKey, buf.Bytes())
	if err != nil {
		return nil, err
	}

	// Parse decrypted body into JSON
	var data map[string]interface{}
	err = sonic.Unmarshal(decryptedBody, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Send a post request without TC Protection
func PostRequestNoTC(url string, body map[string]interface{}) (map[string]interface{}, error) {

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

// Parse encrypted json
func BodyParser(c *fiber.Ctx, data interface{}) error {
	return sonic.Unmarshal(c.Locals("body").([]byte), data)
}

// Return encrypted json
func ReturnJSON(c *fiber.Ctx, data interface{}) error {

	encoded, err := sonic.Marshal(data)
	if err != nil {
		return FailedRequest(c, localization.ErrorServer, err)
	}

	if c.Locals("key") == nil {
		return c.JSON(data)
	}
	encrypted, err := EncryptAES(c.Locals("key").([]byte), encoded)
	if err != nil {
		return FailedRequest(c, localization.ErrorServer, err)
	}

	return c.Send(encrypted)
}

// Get the system uuid set in the environment variables
func GetSystemUUID() uuid.UUID {
	id, err := uuid.Parse(os.Getenv("SYSTEM_UUID"))
	if err != nil {
		panic("Please set the SYSTEM_UUID env property.")
	}

	return id
}
