package integration

import (
	"sync"
	"time"

	"github.com/Liphium/station/pipes"
)

type SessionInformation struct {
	Valid           bool
	Account         string
	PermissionLevel uint

	mutex       *sync.Mutex // A mutex to make other requests wait to not send way too many at the same time
	lastRequest time.Time   // When the last request was sent (mainly for caching)
}

// The response from the server
type sessionResponse struct {
	Success         bool   `json:"success"`
	Account         string `json:"account"`
	PermissionLevel uint   `json:"level"`
}

// Session id -> *Session info (cached for 5 minutes)
var sessionCache = &sync.Map{}

// Validate a session with the id (ideally from the JWT token)
func ValidateSession(node *pipes.LocalNode, session string) *SessionInformation {

	// Check if a session is already in the cache
	var info *SessionInformation
	if obj, exists := sessionCache.Load(session); exists {
		info = obj.(*SessionInformation)
	} else {
		info = &SessionInformation{
			Valid:           false,
			Account:         "",
			PermissionLevel: 0,
			mutex:           &sync.Mutex{},
			lastRequest:     time.Now(),
		}
		sessionCache.Store(session, info)
	}

	// Wait for the mutex to unlock to prevent multiple requests
	info.mutex.Lock()
	defer info.mutex.Unlock()

	// If 5 minutes haven't passed since the last cached request, use the value from there
	if time.Since(info.lastRequest) <= 5*time.Minute {
		return info
	}

	// Make a request to the backend to verify the session
	res, err := PostRequestBackendGeneric[sessionResponse]("/node/get_session", map[string]interface{}{
		"id":      node.ID,
		"token":   node.Token,
		"session": session,
	})
	if err != nil {
		return info
	}

	// Fill the info with the data
	info.Valid = true
	info.Account = res.Account
	info.PermissionLevel = res.PermissionLevel
	info.lastRequest = time.Now()

	return info
}
