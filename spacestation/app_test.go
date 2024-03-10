package main

import (
	"testing"

	"github.com/Liphium/station/spacestation/caching"
)

func TestConcurrency(t *testing.T) {

	// Setup memory
	caching.SetupMemory()
	caching.TestRooms(t)
}
