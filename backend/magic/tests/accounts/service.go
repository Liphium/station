package magic_accounts

import (
	"log"
	"testing"

	"github.com/Liphium/magic/mconfig"
	magic_accounts "github.com/Liphium/station/backend/magic/scripts/accounts"
	magic_util "github.com/Liphium/station/backend/magic/scripts/util"
	"github.com/Liphium/station/backend/service"
	"github.com/Liphium/station/backend/standards"
)

// This test creates an account and makes sure it's returned from the account service.
func MagicService(t *testing.T, p *mconfig.Plan) {
	magic_util.PrepareEnvironment(p)
	magic_util.PrepareDBTest()

	// Test invalid addresses
	t.Run("invalid addresses", func(t *testing.T) {
		if _, err := service.LoadAccount("test1"); err == nil {
			t.Fail()
		}
		if _, err := service.LoadAccount("@@@@@@"); err == nil {
			t.Fail()
		}
		if _, err := service.LoadAccount("test!!!!@"); err == nil {
			t.Fail()
		}
		if _, err := service.LoadAccount("!a!a!"); err == nil {
			t.Fail()
		}
	})

	// Test an account that doesn't exist
	t.Run("invalid account", func(t *testing.T) {
		if _, err := service.LoadAccount(standards.LiphiumAddress("hello")); err == nil {
			t.Fail()
		}
	})

	// Make sure correct info is returned for an actual account
	t.Run("invalid account", func(t *testing.T) {
		id := magic_accounts.TestAccount(p, "test1")

		info, err := service.LoadAccount(standards.LiphiumAddress(id.String()))
		if err != nil {
			log.Fatalln("couldn't load account:", err)
		}

		magic_util.AssertEq(info.Username, "test1")
		magic_util.AssertEq(info.DisplayName, "test1")
		magic_util.AssertEq(info.PublicKey, magic_accounts.DefaultPub)
		magic_util.AssertEq(info.SignatureKey, magic_accounts.DefaultSig)
	})
}
