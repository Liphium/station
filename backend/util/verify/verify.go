package verify

import (
	"errors"
	"sync"
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SessionInformation struct {
	Valid           bool
	Account         string
	PermissionLevel int16

	mutex       *sync.Mutex // A mutex to make other requests wait to not send way too many at the same time
	lastRequest time.Time   // When the last request was sent (mainly for caching)
}

// Session id -> *Session info (cached for 5 minutes)
var sessionCache = &sync.Map{}

func GetSessionInfo(c *fiber.Ctx) (*SessionInformation, error) {

	// Get the session id from the jwt token
	user, valid := c.Locals("user").(*jwt.Token)
	if !valid {
		return nil, errors.New("no jwt token found")
	}
	claims := user.Claims.(jwt.MapClaims)

	// Parse the uuid
	id, err := uuid.Parse(claims["ses"].(string))
	if err != nil {
		return nil, errors.New("couldn't parse UUID: " + err.Error())
	}

	// Check if a session is already in the cache
	var info *SessionInformation
	if obj, exists := sessionCache.Load(id.String()); exists {
		info = obj.(*SessionInformation)
	} else {
		info = &SessionInformation{
			Valid:           false,
			Account:         "",
			PermissionLevel: 0,
			mutex:           &sync.Mutex{},
			lastRequest:     time.Now(),
		}
		sessionCache.Store(id.String(), info)
	}

	// Lock the mutex to prevent multiple database requests for the same thing
	info.mutex.Lock()
	defer info.mutex.Unlock()

	// If 5 minutes haven't passed since the last cached request, use the value from there
	if time.Since(info.lastRequest) <= 5*time.Minute {
		return info, nil
	}

	// Get the session from the database
	var session database.Session
	if err := database.DBConn.Where("id = ?", id).Take(&session).Error; err != nil {
		return nil, err
	}

	// Set the stuff in the session info
	info.Valid = true
	info.Account = session.Account.String()
	info.PermissionLevel = int16(session.PermissionLevel)
	info.lastRequest = time.Now()

	return info, nil
}
