package requests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Liphium/station/backend/settings"
	"github.com/Liphium/station/backend/util"
	"github.com/bytedance/sonic"
)

// Type alias for a map of strings to make it easier to define headers
type Headers = map[string]string

// The protocol used for web requests (either http:// or https://)
var Protocol = "http://"

// The current api version the requests package is using
const ApiVersion = "v1"

// A map used to verify the servers protocol version
var serverCache = &sync.Map{}

// A struct for saving the state of cache, for map above
type cachedServer struct {
	ProtocolVersion uint      // Protocol version of the server
	UpdatedAt       time.Time // When the cache was last refreshed (goes up when no errors)
}

// Duration until server is invalidated in cache
const cacheDuration = 10 * time.Minute

// The response from the /pub endpoint
type publicResponse struct {
	ProtocolVersion uint `json:"protocol_version"`
}

// Send a post request to any server (no generics). Uses Liphium standards.
func PostRequest(server string, path string, body Map) (Map, error) {
	return PostRequestURLGeneric[Map](server+"/"+ApiVersion+path, body)
}

// Send a post request to any server. Uses Liphium standards.
func PostRequestGeneric[T any](server string, path string, body Map) (T, error) {
	var data T

	// Make sure the server follows decentralization requirements
	if !settings.DecentralizationEnabled.ValueOrDefault() {
		return data, fmt.Errorf("decentralization is not allowed")
	}
	if strings.HasPrefix(server, "http://") && !settings.DecentralizationAllowUnsafe.ValueOrDefault() {
		return data, fmt.Errorf("decentralization with unsafe servers not allowed")
	}

	// Make sure there is a protocol specified on the server
	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		server = "https://" + server
	}

	// Make sure the protocol version is correct
	obj, valid := serverCache.Load(server)
	if !valid || time.Since(obj.(cachedServer).UpdatedAt) > cacheDuration {
		res, err := PostRequestURLGeneric[publicResponse](server+"/pub", Map{})
		if err != nil {
			return data, fmt.Errorf("couldn't get protocol version of %s: %s", server, err)
		}
		if res.ProtocolVersion != util.ProtocolVersion {
			return data, fmt.Errorf("protocol versions are incompatible: %d (current) vs. %d (target server)", util.ProtocolVersion, res.ProtocolVersion)
		}

		// Add the up to date information
		obj = cachedServer{
			ProtocolVersion: res.ProtocolVersion,
			UpdatedAt:       time.Now(),
		}
		serverCache.Store(server, obj)
	}

	// Send the actual request
	return PostRequestURLGeneric[T](server+"/"+ApiVersion+path, body)
}

// Send a post request to any URL
func PostRequestURL(url string, body Map) (Map, error) {
	return PostRequestURLGenericWithHeaders[Map](url, body, Headers{})
}

// Send a post request to any URL
func PostRequestURLGeneric[T any](url string, body Map) (T, error) {
	return PostRequestURLGenericWithHeaders[T](url, body, Headers{})
}

// Send a post request to any URL with headers attached
func PostRequestURLGenericWithHeaders[T any](url string, body Map, headers Headers) (T, error) {

	// Declared here so it can be returned as nil before it's actually used
	var data T

	// Encode body to JSON
	byteBody, err := sonic.Marshal(body)
	if err != nil {
		return data, err
	}

	// Set headers
	reqHeaders := http.Header{}
	reqHeaders.Set("Content-Type", "application/json")
	for key, value := range headers {
		reqHeaders.Set(key, value)
	}

	// Send the request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(byteBody))
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

// A useful helper struct for the normal response you get from the server (use with generics).
type NormalResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
