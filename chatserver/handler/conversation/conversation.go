package conversation

import (
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/handler/conversation/space"
	message_routes "github.com/Liphium/station/chatserver/routes/conversations/message"
	"github.com/Liphium/station/pipeshandler/wshandler"
)

func SetupActions() {
	space.SetupActions()

	wshandler.RegisterHandler(caching.Node, "conv_sub", subscribe)

	// Setup messages queue
	setupMessageQueue()
}

const messageProcessorAmount = 3

type TokenTask struct {
	Adapter      string
	Conversation string
	Date         int64 // Unix timestamp of last fetch
}

var newTaskChan = make(chan TokenTask)

func setupMessageQueue() {
	for i := 0; i < messageProcessorAmount; i++ {
		go func() {
			for {

				// Wait for a new task
				task := <-newTaskChan

				// Get all messages
				var messages []conversations.Message
				if database.DBConn.Where("conversation = ? AND creation > ?", task.Conversation, task.Date).Find(&messages).Error != nil {
					continue
				}

				// Send messages to the adapter
				for _, message := range messages {
					caching.Node.SendClient(task.Adapter, message_routes.MessageEvent(message))
					time.Sleep(3 * time.Millisecond) // Give TCP some time to send the message
				}
			}
		}()
	}
}

func AddConversationToken(task TokenTask) {
	newTaskChan <- task
}
