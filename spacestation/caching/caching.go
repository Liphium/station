package caching

import (
	"github.com/Liphium/station/spacestation/util"
)

func SetupMemory() {
	setupRoomsCache()
	setupRoomConnectionsCache()
	setupConnectionsCache()
	setupSessionsCache()
	setupTablesCache()
}

func CloseCaches() {
	util.Log.Println("Closing caches...")
	roomsCache.Close()
	roomConnectionsCache.Close()
	connectionsCache.Close()
	sessionsCache.Close()
	tablesCache.Close()
}
