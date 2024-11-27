package caching

import (
	"sync"

	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
)

type SyncData struct {
	TokenID      string
	Conversation string
	Since        int64
}

// The stuff needed to manage threads
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

	// Make space for a new thread after work is finished
	defer func() {
		<-currentChan
	}()

	// Sync all messages until there are none left
	finished := false
	for !finished {

		// Get the messages
		var messages []conversations.Message
		if err := database.DBConn.Where("creation >= ?", data.Since).Order("creation ASC").Limit(30).Find(&messages).Error; err != nil {
			return err
		}

		// Send an event to the client containing the messages
		if err := CSNode.Pipe(pipes.ProtocolWS, pipes.Message{
			Local:   true,
			Channel: pipes.BroadcastChannel([]string{"s-" + data.TokenID}),
			Event: pipes.Event{
				Name: "conv_msg_mp",
				Data: map[string]interface{}{
					"cv":   data.Conversation,
					"msgs": messages,
				},
			},
		}); err != nil {
			return err
		}

		// Update state for the next iteration
		finished = len(messages) < 30
		data.Since = messages[len(messages)-1].Creation
	}

	return nil
}
