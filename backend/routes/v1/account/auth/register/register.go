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
	router.Post("/start", startRegister)
	router.Post("/email", checkEmail)
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
