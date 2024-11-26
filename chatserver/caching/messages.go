package caching

import (
	"sync"

	"github.com/Liphium/station/main/integration"
)

var mutexLock = &sync.Mutex{}
var mutexCache = []*sync.Mutex{}

func AddMessageToCache() error {
	mutexLock.Lock()

	_, err := integration.GetIntSetting(CSNode, integration.SettingChatMessagePullThreads)
	if err != nil {
		return err
	}

	mutexLock.Unlock()

	return nil
}
