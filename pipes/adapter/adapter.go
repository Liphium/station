package adapter

import (
	"log"
	"sync"

	"github.com/Liphium/station/pipes"
	"github.com/dgraph-io/ristretto"
)

type Adapter struct {
	ID    string      // Identifier of the client
	Mutex *sync.Mutex // Mutex to prevent concurrent sending (WHY DO I NEED TO DO THIS??)
	Data  interface{} // Custom data (not required)

	// Functions
	Receive func(*Context) error
}

type Context struct {
	Event   *pipes.Event
	Message []byte
	Adapter *Adapter
}

var websocketCache *ristretto.Cache
var udpCache *ristretto.Cache

func SetupCaching() {
	var err error
	websocketCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,
	})

	if err != nil {
		panic(err)
	}

	udpCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,
	})

	if err != nil {
		panic(err)
	}
}

// Register a new adapter for websocket/sl (all safe protocols)
func AdaptWS(adapter Adapter) {
	if adapter.Mutex == nil {
		adapter.Mutex = &sync.Mutex{}
	}

	if websocketCache == nil {
		panic("Please call adapter.SetupCaching() before using the adapter package")
	}

	_, ok := websocketCache.Get(adapter.ID)
	if ok {
		websocketCache.Del(adapter.ID)
		log.Printf("[ws] Replacing adapter for target %s \n", adapter.ID)
	}

	websocketCache.Set(adapter.ID, &adapter, 1)
}

// Register a new adapter for UDP
func AdaptUDP(adapter Adapter) {
	if adapter.Mutex == nil {
		adapter.Mutex = &sync.Mutex{}
	}

	if websocketCache == nil {
		panic("Please call adapter.SetupCaching() before using the adapter package")
	}

	_, ok := udpCache.Get(adapter.ID)
	if ok {
		udpCache.Del(adapter.ID)
		log.Printf("[udp] Replacing adapter for target %s \n", adapter.ID)
	}

	udpCache.Set(adapter.ID, &adapter, 1)
}

// Remove a websocket/sl adapter
func RemoveWS(ID string) {
	websocketCache.Del(ID)
}

// Remove a UDP adapter
func RemoveUDP(ID string) {
	udpCache.Del(ID)
}

// Handles receiving messages from the target and passes them to the adapter
func ReceiveWeb(ID string, event pipes.Event, msg []byte) {

	obj, ok := websocketCache.Get(ID)
	if !ok {
		return
	}

	adapter := obj.(*Adapter)
	adapter.Mutex.Lock()
	err := adapter.Receive(&Context{
		Event:   &event,
		Message: msg,
		Adapter: adapter,
	})
	adapter.Mutex.Unlock()

	if err != nil {
		log.Printf("[ws] Error receiving message from target %s: %s \n", ID, err)
	}
}

// Handles receiving messages from the target and passes them to the adapter
func ReceiveUDP(ID string, event pipes.Event, msg []byte) {

	obj, ok := udpCache.Get(ID)
	if !ok {
		return
	}

	adapter := obj.(*Adapter)
	adapter.Mutex.Lock()
	err := adapter.Receive(&Context{
		Event:   &event,
		Message: msg,
		Adapter: adapter,
	})
	adapter.Mutex.Unlock()

	if err != nil {
		log.Printf("[udp] Error receiving message from target %s: %s \n", ID, err)
	}
}
