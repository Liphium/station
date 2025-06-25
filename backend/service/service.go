package service

import (
	"fmt"
	"os"

	"github.com/Liphium/station/neogate"
)

var Instance *neogate.Instance

// Get the Liphium address for this server and an account id.
func LiphiumAddress(accountId string) string {
	protocol := os.Getenv("PROTOOCOL")
	if protocol == "http://" {
		return fmt.Sprintf("%s@%s%s", accountId, protocol, os.Getenv("BASE_PATH"))
	}
	return fmt.Sprintf("%s@%s", accountId, os.Getenv("BASE_PATH"))
}
