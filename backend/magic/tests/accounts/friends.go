package magic_accounts

import (
	"testing"

	"github.com/Liphium/magic/mconfig"
	magic_accounts "github.com/Liphium/station/backend/magic/scripts/accounts"
	magic_util "github.com/Liphium/station/backend/magic/scripts/util"
	magic_gate "github.com/Liphium/station/backend/magic/tests/gate"
	friends2_routes "github.com/Liphium/station/backend/routes/v1/accounts/friends2"
	"github.com/Liphium/station/backend/service"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util/requests"
)

// Do not call this function anything with Test, it will cause errors
func MagicFriendAdding(t *testing.T, p *mconfig.Plan) {
	magic_util.PrepareEnvironment(p)
	magic_util.PrepareDBTest()

	// Prepare
	test1 := magic_accounts.TestAccount(p, "test1")
	test1LPH := standards.LiphiumAddress(test1.ID.String())
	_, tk := magic_accounts.GetTestToken(p, "test1")
	test2 := magic_accounts.TestAccount(p, "test2")
	test2LPH := standards.LiphiumAddress(test2.ID.String())
	_, tk2 := magic_accounts.GetTestToken(p, "test2")

	test1Gate := magic_gate.NewGateConnection(t, tk)
	defer test1Gate.Close(t)
	test2Gate := magic_gate.NewGateConnection(t, tk2)
	defer test2Gate.Close(t)

	t.Run("friend request sending", func(t *testing.T) {
		res, err := requests.PostRequestAuthURL(requests.CurrentPath("/a/accounts/friends/add"), tk, friends2_routes.FriendAddRequest{
			Id:         test2LPH,
			ProfileKey: "prof",
		})
		magic_util.AssertEq(t, err, nil)
		magic_util.AssertEq(t, requests.ValueOr(res, "success", false), true)

		frEv := test2Gate.ReadEvent(t)
		test1Info, _, err := service.LoadAccount(test1LPH)
		magic_util.AccountServiceError(t, err)
		magic_util.AssertDeepEq(t, frEv, friends2_routes.FriendEventFromAccountInfo(true, test1Info, "prof"))
	})

	t.Run("friend request accepting", func(t *testing.T) {
		res, err := requests.PostRequestAuthURL(requests.CurrentPath("/a/accounts/friends/add"), tk2, friends2_routes.FriendAddRequest{
			Id:         test1LPH,
			ProfileKey: "prof",
		})
		magic_util.AssertEq(t, err, nil)
		magic_util.AssertEq(t, requests.ValueOr(res, "success", false), true)

		// Event for test1
		frEv := test1Gate.ReadEvent(t)
		test2Info, _, err := service.LoadAccount(test2LPH)
		magic_util.AccountServiceError(t, err)
		magic_util.AssertDeepEq(t, frEv, friends2_routes.FriendEventFromAccountInfo(false, test2Info, "prof"))

		// Event for test2
		frEv = test2Gate.ReadEvent(t)
		test1Info, _, err := service.LoadAccount(test1LPH)
		magic_util.AccountServiceError(t, err)
		magic_util.AssertDeepEq(t, frEv, friends2_routes.FriendEventFromAccountInfo(false, test1Info, "prof"))
	})
}
