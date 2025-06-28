package magic_util

import (
	"os"

	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/station/backend/database"
	backend_starter "github.com/Liphium/station/backend/starter"
)

// Add all environments from the plan to the actual environment.
func PrepareEnvironment(p *mconfig.Plan) {
	for k, v := range p.Environment {
		os.Setenv(k, v)
	}
}

func PrepareDBTest() {
	database.Connect()
	backend_starter.CreateDefaultObjects()
}
