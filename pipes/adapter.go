package pipes

import (
	"sync"

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
	Event   *Event
	Message []byte
	Adapter *Adapter
}

func (node *LocalNode) setupCaching() {
	var err error
	node.websocketCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,
	})

	if err != nil {
		panic(err)
	}
}

// Register a new adapter for websocket/sl (all safe protocols)
func (node *LocalNode) AdaptWS(adapter Adapter) {
	if adapter.Mutex == nil {
		adapter.Mutex = &sync.Mutex{}
	}

	if node.websocketCache == nil {
		panic("Please call adapter.SetupCaching() before using the adapter package")
	}

	_, ok := node.websocketCache.Get(adapter.ID)
	if ok {
		node.websocketCache.Del(adapter.ID)
		Log.Printf("[ws] Replacing adapter for target %s \n", adapter.ID)
	}

	node.websocketCache.Set(adapter.ID, &adapter, 1)
}

// Remove a websocket/sl adapter
func (node *LocalNode) RemoveAdapterWS(ID string) {
	node.websocketCache.Del(ID)
}

// Handles receiving messages from the target and passes them to the adapter
func (node *LocalNode) AdapterReceiveWeb(ID string, event Event, msg []byte) {

	obj, ok := node.websocketCache.Get(ID)
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
		Log.Printf("[ws] Error receiving message from target %s: %s \n", ID, err)
	}
}
