package integration

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
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

// Grab public key from the server.
func grabServerPublicKey() error {
	res, err := http.Post(BasePath+"/"+ApiVersion+"/pub", "application/json", nil)
	if err != nil {
		return err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, res.Body)

	if err != nil {
		return err
	}

	var data map[string]interface{}
	err = sonic.Unmarshal([]byte(buf.String()), &data)
	if err != nil {
		return err
	}

	ServerPublicKey, err = UnpackageRSAPublicKey(data["pub"].(string))
	if err != nil {
		return err
	}

	return res.Body.Close()
}

var Protocol = "http://"
var BasePath = "http://localhost:3000"
var ServerPublicKey *rsa.PublicKey // Public key from the backend server

// * Important
const ApiVersion = "v1"

// Send a post request (without TC encryption and custom URL)
func PostRequestNoTC(url string, body map[string]interface{}) (map[string]interface{}, error) {

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

// Send a post request (with TC protection encryption)
func PostRequest(url string, body map[string]interface{}) (map[string]interface{}, error) {

	if ServerPublicKey == nil {
		panic("Server public key not set.")
	}

	byteBody, err := sonic.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Compute the auth tag
	aesKey, err := NewAESKey()
	if err != nil {
		return nil, err
	}
	authTag, err := EncryptRSA(ServerPublicKey, aesKey)
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

	// Send the request
	req, err := http.NewRequest(http.MethodPost, BasePath+"/"+ApiVersion+url, reader)
	if err != nil {
		return nil, err
	}
	req.Header = reqHeaders

	res, err := http.DefaultClient.Do(req)
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

type NormalResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
