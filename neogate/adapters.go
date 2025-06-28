package neogate

import (
	"errors"
	"sync"

	"github.com/bytedance/sonic"
)

type AdapterFunc = func(*AdapterContext) error

type Adapter struct {
	ID    string      // Identifier of the client
	Mutex *sync.Mutex // Mutex to prevent concurrent exceptions (can happen with connections, better handle this on the neogate level)

	// Functions
	OnEvent AdapterFunc
	OnError func(error)
}

type AdapterContext struct {
	Event   *Event
	Message []byte
	Adapter *Adapter
}

type Event struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}

type CreateAction struct {
	ID      string      // Id of the adapter
	OnEvent AdapterFunc // Function that handles events received by the adapter
	OnError func(error) // Function that handles errors encountered by the adapter
}

// Register a new adapter for websocket/sl (all safe protocols)
func (instance *Instance) Adapt(createAction CreateAction) {
	_, ok := instance.adapters.Load(createAction.ID)
	if ok {
		instance.adapters.Delete(createAction.ID)
		Log.Printf("Replacing adapter for target %s \n", createAction.ID)
	}

	instance.adapters.Store(createAction.ID, &Adapter{
		ID:      createAction.ID,
		Mutex:   &sync.Mutex{},
		OnEvent: createAction.OnEvent,
		OnError: createAction.OnError,
	})
}

// Remove an adapter from the instance
func (instance *Instance) RemoveAdapter(ID string) {
	instance.adapters.Delete(ID)
}

// Handles receiving messages from the target and passes them to the adapter
func (instance *Instance) AdapterReceive(ID string, event Event, msg []byte) error {

	obj, ok := instance.adapters.Load(ID)
	if !ok {
		return errors.New("adapter not found")
	}
	adapter := obj.(*Adapter)

	adapter.Mutex.Lock()
	defer adapter.Mutex.Unlock()

	err := adapter.OnEvent(&AdapterContext{
		Event:   &event,
		Message: msg,
		Adapter: adapter,
	})

	// Tell the adapter there was an error
	if err != nil {
		adapter.OnError(err)
		Log.Printf("[ws] Error receiving message from target %s: %s \n", ID, err)
	}
	return err
}

// Send an event to all adapters
func (instance *Instance) Send(adapters []string, event Event) error {
	msg, err := sonic.Marshal(event)
	if err != nil {
		return err
	}

	for _, adapter := range adapters {
		instance.AdapterReceive(adapter, event, msg)
	}
	return nil
}

// Sends an event to the account.
//
// Only returns errors for encoding, not retrieval (cause adapters handle that themselves).
func (instance *Instance) SendOne(adapter string, event Event) error {
	return instance.Send([]string{adapter}, event)
}
