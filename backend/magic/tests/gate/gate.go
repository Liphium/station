package magic_gate

import (
	"testing"

	"github.com/Liphium/magic/mconfig"
	magic_accounts "github.com/Liphium/station/backend/magic/scripts/accounts"
	magic_util "github.com/Liphium/station/backend/magic/scripts/util"
)

func MagicGateConnection(p *mconfig.Plan, t *testing.T) {
	magic_util.PrepareEnvironment(p)
	magic_util.PrepareDBTest()

	t.Run("basic gate flow", func(t *testing.T) {
		magic_accounts.TestAccount(p, "test1")
		_, tk := magic_accounts.GetTestToken(p, "test1")

		// Validate connections work
		conn := magic_util.NewGateConnection(t, tk)
		defer conn.Close(t)
	})
}
