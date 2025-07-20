package magic_events

import (
	"fmt"

	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/station/backend/database"
	magic_accounts "github.com/Liphium/station/backend/magic/scripts/accounts"
	magic_util "github.com/Liphium/station/backend/magic/scripts/util"
)

// Create a test token for the account with that username
func RunEvents(p *mconfig.Plan, username string) {
	magic_util.PrepareEnvironment(p)
	database.Connect()

	_, tk := magic_accounts.GetTestToken(p, username)
	gate := magic_util.NewGateConnection(nil, tk)

	for {
		ev := gate.ReadEvent(nil)
		fmt.Println(ev)
	}
}
