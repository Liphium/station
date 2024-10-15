package verify

import (
	"errors"
	"sync"
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SessionInformation struct {
	valid           bool
	account         string
	session         string
	permissionLevel int16

	mutex       *sync.Mutex // A mutex to make other requests wait to not send way too many at the same time
	lastRequest time.Time   // When the last request was sent (mainly for caching)
}

func (i *SessionInformation) IsValid() bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.valid
}

func (i *SessionInformation) GetAccount() string {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.account
}

func (i *SessionInformation) GetAccountUUID() (uuid.UUID, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	// Convert the account id to a UUID
	return uuid.Parse(i.account)
}

func (i *SessionInformation) GetSession() string {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.session
}

func (i *SessionInformation) GetSessionUUID() (uuid.UUID, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	// Convert the session id to a UUID
	return uuid.Parse(i.session)
}

func (i *SessionInformation) GetPermissionLevel() int16 {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.permissionLevel
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
			valid:           false,
			account:         "",
			session:         id.String(),
			permissionLevel: 0,
			mutex:           &sync.Mutex{},
			lastRequest:     time.UnixMilli(0),
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
	info.valid = true
	info.account = session.Account.String()
	info.session = session.ID.String()
	info.permissionLevel = int16(session.PermissionLevel)
	info.lastRequest = time.Now()

	return info, nil
}

// A middleware that checks the jwt token and then also loads all of the session information into the locals
func AuthMiddleware() func(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS512,
			Key:    []byte(util.JWT_SECRET),
		},

		// Checks if the token is expired
		SuccessHandler: func(c *fiber.Ctx) error {

			// Make sure the token isn't expired
			if util.IsExpired(c) {
				return util.InvalidRequest(c)
			}

			// Verify the session
			info, err := GetSessionInfo(c)
			if err != nil {
				if util.Testing {
					util.Log.Println("invalid session info: " + err.Error())
				}
				return util.InvalidRequest(c)
			}

			// Make sure the session is valid
			if !info.IsValid() {
				if util.Testing {
					util.Log.Println("the session is invalid")
				}
				return util.InvalidRequest(c)
			}

			// Add the info to the locals
			c.Locals("info", info)

			return c.Next()
		},

		// Error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			util.Log.Println(err.Error())

			// Return error message
			return c.SendStatus(401)
		},
	})
}

// Get the session information from the locals (just makes the process a little easier)
func InfoLocals(c *fiber.Ctx) *SessionInformation {
	return c.Locals("info").(*SessionInformation)
}
