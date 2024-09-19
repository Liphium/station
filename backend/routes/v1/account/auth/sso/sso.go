package sso_routes

import (
	"log"
	"os"
	"strings"

	"github.com/Liphium/station/backend/kv"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
)

// I'd like to say this here because it's the easiest place to say it.
// Most of this code was heavily inspired by https://github.com/Shareed2k/goth_fiber.
// Thanks for making this amazing thing. I'm just not using your library cause my
// backend is very weird and needs a custom implementation for this use-case but
// because I still wanted to mention you, this is here. If you are not credited on the
// credits page on the Liphium website when it's finally done, please let me know
// and we shall change that.

var openIdName = ""
var Enabled = false

func Unencrypted(router fiber.Router) {

	if os.Getenv("SSO_ENABLED") != "true" {
		return
	}

	// Setup the open id provider
	id := os.Getenv("SSO_CLIENT_ID")
	secret := os.Getenv("SSO_CLIENT_SECRET")
	url := os.Getenv("SSO_CONFIG")

	scopes := []string{"email"}
	if os.Getenv("SSO_SCOPES") != "" {
		scopes = strings.Split(os.Getenv("SSO_SCOPES"), " ")
	}

	// Add the open id provider to goth
	openIdProvider, err := openidConnect.New(id, secret, os.Getenv("PROTOCOL")+os.Getenv("BASE_PATH")+"/v1/account/auth/sso/callback", url, scopes...)
	if err != nil {
		panic(err)
	}
	if openIdProvider != nil {
		openIdName = openIdProvider.Name()
		goth.UseProviders(openIdProvider)
	} else {
		panic("open id provider couldn't be set up")
	}

	// Set it to enabled
	Enabled = true

	log.Println("SSO is enabled")

	// Register the callback endpoint
	router.Get("/callback", callback)
}

func Unauthorized(router fiber.Router) {

	if os.Getenv("SSO_ENABLED") != "true" {
		return
	}

	// Register all the endpoints for SSR
	router.Post("/form", getSSOForm)
	router.Post("/check", checkSSO)
}

// Implementation of goth.Params copied from goth_fiber (look at notice above)
type Params struct {
	ctx *fiber.Ctx
}

func (p *Params) Get(key string) string {
	return p.ctx.Query(key)
}

type SSOState struct {
	// Stuff for the redirect itself
	State   string // The state string provided to the open id provider
	Session string // The goth (auth library) session
	Locale  string // The locale of the thing (can't be fetched in the callback endpoint)

	// Stuff for the check endpoint
	UserID string // The user id from the SSO provider
	Email  string // The email from the SSO provider
	Done   bool   // Check if the thing is done
}

const ssoTokenPrefix = "sso_"

// Generate the token required for using SSO with SSR
func GenerateSSOToken(c *fiber.Ctx) string {

	// Generate a unique token
	token := auth.GenerateToken(50)
	for _, valid := kv.Get(ssoTokenPrefix + token); valid; {
		token = auth.GenerateToken(50)
	}

	// Store it with a completed bool
	kv.Store(ssoTokenPrefix+token, &SSOState{
		State:  token,
		Locale: localization.Locale(c),
		Done:   false,
	})
	return token
}

// Check a token and see if SSO was completed
func checkToken(token string) (*SSOState, localization.Translations) {

	// Get the token
	obj, valid := kv.Get(ssoTokenPrefix + token)
	if !valid {
		return nil, localization.ErrorInvalidRequest
	}

	// Return whether SSO was completed
	return obj.(*SSOState), nil
}
