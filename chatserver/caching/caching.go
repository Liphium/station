package caching

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// This just needs to be kept somewhere to avoid import cycles
var CSInstance *pipeshandler.Instance
var CSNode *pipes.LocalNode
