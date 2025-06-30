package magic_util

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/station/backend/database"
	backend_starter "github.com/Liphium/station/backend/starter"
	"github.com/Liphium/station/backend/util/requests"
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

// Only for localhost
func GatewayURL() string {
	return strings.ReplaceAll(requests.CurrentPath("/gate/connect"), "http://", "ws://")
}

// Marshaling for json without error handling cause testing (will just panic instead)
func Marshal(data interface{}) []byte {
	done, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("couldn't marshal:", err)
	}
	return done
}

// Unmarshalling for json without error handling cause testing (will just panic instead)
func Unmarshal(data []byte, v any) {
	if err := json.Unmarshal(data, v); err != nil {
		log.Fatalln("couldn't unmarshal:", err)
	}
}
