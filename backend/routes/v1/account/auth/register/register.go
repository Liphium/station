package register_routes

import (
	"sync"
	"time"

	"github.com/Liphium/station/backend/kv"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(router fiber.Router) {
	router.Post("/email", checkEmail)          // Step 0: Check the email and redirect to save the token
	router.Post("/start", startRegister)       // Step 1: Render the invite form
	router.Post("/invite", checkInvite)        // Step 2: Check the invite code
	router.Post("/email_code", checkEmailCode) // Step 3: Check the email code
	router.Post("/resend_email", resendEmail)  // Step 3: Resend email endpoint
	router.Post("/username", checkUsername)    // Step 4: Check the username
	router.Post("/password", checkPassword)    // Step 5: Check the password & return tokens
}

type RegisterState struct {
	Step  uint        // The step of the register process
	Mutex *sync.Mutex // To prevent frequent requests killing the server

	// Simple data to enter
	Invite      string // The invite used to create the account
	Username    string // The username of the account
	DisplayName string // The display name of the account

	// Email verification stuff
	Email     string    // The email of the account
	EmailCode string    // The code required for verification
	LastEmail time.Time // The last time an email was sent

	// Rate limiting for whatever endpoint the guy is on right now
	AttemptCount uint      // The count of attempts
	LastAttempt  time.Time // The last attempt to get in
}

const registerTokenPrefix = "register_"

// Generate the token required for logging in with SSR
func GenerateRegisterToken(email string) string {

	// Generate a unique token
	token := auth.GenerateToken(50)
	for _, valid := kv.Get(registerTokenPrefix + token); valid; {
		token = auth.GenerateToken(50)
	}

	// Store it as a login token in the kv store
	kv.Store(registerTokenPrefix+token, &RegisterState{
		Step:  1,
		Mutex: &sync.Mutex{},
		Email: email,
	})
	return token
}

// Make sure the user has the correct token and permission for the current endpoint
func validateToken(token string, step uint) (*RegisterState, localization.Translations) {

	// Get the token from the key-value store
	obj, valid := kv.Get(registerTokenPrefix + token)
	if !valid {
		return nil, localization.ErrorInvalidRequest
	}

	// Check if the user can access that endpoint
	state := obj.(*RegisterState)
	if state.Step != step {
		return nil, localization.ErrorInvalidRequest
	}

	return state, nil
}

// Upgrade a token to a higher step
func upgradeToken(token string, step uint) localization.Translations {

	// Get the token from the key-value store
	obj, valid := kv.Get(registerTokenPrefix + token)
	if !valid {
		return localization.ErrorInvalidRequest
	}

	// Check if the user can access that endpoint
	state := obj.(*RegisterState)

	// Lock the mutex to prevent modification and update
	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	state.Step = step

	// Also reset rate limiting
	state.AttemptCount = 0
	state.LastAttempt = time.Now()

	return nil
}

// Make sure the user doesn't go over the rate limit
func ratelimitHandler(state *RegisterState, maxAttempts uint, cooldown time.Duration) bool {

	// Prevent concurrent reads/writes
	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	// Check if there have been too many attempts
	if state.AttemptCount > maxAttempts {

		// Check if the rate limit can already be reset
		if time.Since(state.LastAttempt) > cooldown {
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
