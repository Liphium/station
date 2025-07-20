package magic_accounts2

import (
	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/station/backend/database"
	magic_accounts "github.com/Liphium/station/backend/magic/scripts/accounts"
	magic_util "github.com/Liphium/station/backend/magic/scripts/util"
	friends2_routes "github.com/Liphium/station/backend/routes/v1/accounts/friends2"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util/requests"
)

func RunFriendTest(p *mconfig.Plan) {
	magic_util.PrepareEnvironment(p)
	database.Connect()

	// Prepare
	_, tk := magic_accounts.GetTestToken(p, "test1")
	test2LPH := standards.LiphiumAddress("7b2687b9-688c-4970-957d-9825f75e5de5")

	requests.PostRequestAuthURL(requests.CurrentPath("/a/accounts/friends/add"), tk, friends2_routes.FriendAddRequest{
		Id:         test2LPH,
		ProfileKey: "prof",
	})
}
