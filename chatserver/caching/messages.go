package caching

import (
	"sync"

	"github.com/Liphium/station/main/integration"
)

type SyncData struct {
	Conversation string
	Since        int64
}

var mutexLock = &sync.Mutex{}
var bufferedChannel chan struct{}

// Always start this method in a new goroutine (it assumes you do)
func AddSyncToQueue(data SyncData) error {

	// Get the amount of message pull threads
	threadsSet, err := integration.GetIntSetting(CSNode, integration.SettingChatMessagePullThreads)
	if err != nil {
		return err
	}
	threads := int(threadsSet)

	// Make sure no concurrent writes happen
	mutexLock.Lock()

	// Update the size of the channel if the setting changed
	if bufferedChannel == nil || cap(bufferedChannel) != threads {
		bufferedChannel = make(chan struct{}, threads)
	}

	// Unlock the thing to make sure future goroutines can use it
	mutexLock.Unlock()

	// Add something to the channel (will block until space is available / threads are available)
	currentChan := bufferedChannel
	currentChan <- struct{}{}

	// Make space for a new thread
	<-currentChan

	return nil
}
