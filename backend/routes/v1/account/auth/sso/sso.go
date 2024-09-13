package sso_routes

import (
	"os"

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

func Unencrypted(router fiber.Router) {

	if os.Getenv("SSO_ENABLED") != "true" {
		return
	}

	// Setup the open id provider
	id := os.Getenv("SSO_CLIENT_ID")
	secret := os.Getenv("SSO_CLIENT_SECRET")
	url := os.Getenv("SSO_CONFIG")

	// Add the open id provider to goth
	openIdProvider, err := openidConnect.New(id, secret, os.Getenv("PROTOCOL")+os.Getenv("BASE_PATH")+"/v1/account/auth/sso/callback", url)
	if err != nil {
		panic(err)
	}
	if openIdProvider != nil {
		openIdName = openIdProvider.Name()
		goth.UseProviders(openIdProvider)
	} else {
		panic("open id provider couldn't be set up")
	}
}
