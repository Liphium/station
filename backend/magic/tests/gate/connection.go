package magic_gate

import (
	"testing"

	"github.com/Liphium/magic/mconfig"
	magic_util "github.com/Liphium/station/backend/magic/scripts/util"
)

// Do not call this function anything with Test, it will cause errors
func MagicConnection(t *testing.T, p *mconfig.Plan) {
	magic_util.PrepareEnvironment(p)
	magic_util.PrepareDBTest()

	t.Run("basic gate flow", func(t *testing.T) {

	})

}
