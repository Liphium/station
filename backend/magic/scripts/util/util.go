package magic_util

import (
	"os"

	"github.com/Liphium/magic/mconfig"
)

// Add all environments from the plan to the actual environment.
func PrepareEnvironment(p *mconfig.Plan) {
	for k, v := range p.Environment {
		os.Setenv(k, v)
	}
}
