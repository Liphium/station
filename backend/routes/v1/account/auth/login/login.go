package login_routes

import (
	"sync"
	"time"

	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/kv"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Register all the routes related to logging in with SSR
func Unauthorized(router fiber.Router) {
	router.Post("/password", checkPassword)
	router.Post("/start", startLogin)
}

type LoginState struct {
	LoginStep       uint        // The step of the login process
	Account         uuid.UUID   // The account id for the account
	PermissionLevel uint        // The permission level of the account
	Mutex           *sync.Mutex // To prevent the rate limit triggering concurrent reads
	AttemptCount    uint        // The count of attempts
	LastAttempt     time.Time   // The last attempt to get in
}

const loginTokenPrefix = "login_"

// Generate the token required for logging in with SSR
func GenerateLoginToken(acc account.Account) string {

	// Generate a unique token
	token := auth.GenerateToken(50)
	for _, valid := kv.Get(token); valid; {
		token = auth.GenerateToken(50)
	}

	// Store it as a login token in the kv store
	kv.Store(loginTokenPrefix+token, &LoginState{
		LoginStep:       1,
		Account:         acc.ID,
		PermissionLevel: acc.Rank.Level,
		Mutex:           &sync.Mutex{},
		AttemptCount:    0,
		LastAttempt:     time.Now(),
	})
	return token
}

// Make sure the user has the correct token and permission for the current endpoint
func testTokenAndRatelimit(token string, step uint) localization.Translations {

	// Get the token from the key-value store
	obj, valid := kv.Get(loginTokenPrefix + token)
	if !valid {
		return localization.ErrorInvalidRequest
	}

	// Check if the user can access that endpoint
	state := obj.(*LoginState)
	if state.LoginStep != step {
		return localization.ErrorInvalidRequest
	}
	if !ratelimitHandler(state) {
		return localization.ErrorAuthRatelimit
	}

	return nil
}

// Make sure the user doesn't go over the rate limit
func ratelimitHandler(state *LoginState) bool {

	// Prevent concurrent reads/writes
	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	// Check if there have been too many attempts
	if state.AttemptCount > 3 {

		// Check if the rate limit can already be reset
		if time.Since(state.LastAttempt) > time.Second*3 {
			state.AttemptCount = 0
			state.LastAttempt = time.Now()
			return true
		}

		return false
	}

	// Update the rate limit data
	state.AttemptCount++
	state.LastAttempt = time.Now()

	return true
}
