package conversation

import (
	"github.com/Liphium/station/pipes/receive/processors"
)

func SetupProcessors() {
	processors.Processors["conv_open:l"] = open
}
