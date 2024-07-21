package caching

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/util"
)

var SSInstance *pipeshandler.Instance
var SSNode *pipes.LocalNode

func SetupMemory() {
	setupRoomsCache()
	setupRoomConnectionsCache()
	setupConnectionsCache()
	setupSessionsCache()
}

func CloseCaches() {
	util.Log.Println("Closing caches...")
	roomsCache.Close()
	roomConnectionsCache.Close()
	connectionsCache.Close()
	sessionsCache.Close()
}
