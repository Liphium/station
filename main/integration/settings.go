package integration

import (
	"sync"
	"time"

	"github.com/Liphium/station/pipes"
)

// Predefined settings
const SettingDecentralizationEnabled = "decentralization.enabled"
const SettingDecentralizationUnsafeAllowed = "decentralization.allow_unsafe"
const SettingChatMessagePullThreads = "chat.message_pull_threads"

type boolSettingCacheData struct {
	Value       bool      // The value of the setting
	LastRequest time.Time // The time of the last request
}

// Name -> Cache data
var boolSettingCache = &sync.Map{}

type getSettingBoolResponse struct {
	Success bool `json:"success"`
	Value   bool `json:"value"`
}

// Get any setting on the backend and parse it
func GetBoolSetting(node *pipes.LocalNode, name string) (bool, error) {

	// Check if there is already a value in cache
	if obj, valid := boolSettingCache.Load(name); valid {
		data := obj.(boolSettingCacheData)

		// If 5 minutes haven't passed since the last cached request, use the value from there
		if time.Since(data.LastRequest) <= 5*time.Minute {
			return data.Value, nil
		} else {
			// Otherwise, delete the data in the cache and do another request
			boolSettingCache.Delete(name)
		}
	}

	// Get the thing from the backend instead if not cached
	res, err := PostRequestBackendGeneric[getSettingBoolResponse]("/node/get_bool_setting", map[string]interface{}{
		"id":      node.ID,
		"token":   node.Token,
		"setting": name,
	})
	if err != nil {
		return false, err
	}

	// Insert it into the cache
	boolSettingCache.Store(name, boolSettingCacheData{
		Value:       res.Value,
		LastRequest: time.Now(),
	})

	// Parse the result
	return res.Value, err
}

type intSettingCacheData struct {
	Value       int64     // The value of the setting
	LastRequest time.Time // The time of the last request
}

// Name -> Cache data
var intSettingCache = &sync.Map{}

type getSettingIntResponse struct {
	Success bool  `json:"success"`
	Value   int64 `json:"value"`
}

// Get any setting on the backend and parse it
func GetIntSetting(node *pipes.LocalNode, name string) (int64, error) {

	// Check if there is already a value in cache
	if obj, valid := intSettingCache.Load(name); valid {
		data := obj.(intSettingCacheData)

		// If 5 minutes haven't passed since the last cached request, use the value from there
		if time.Since(data.LastRequest) <= 5*time.Minute {
			return data.Value, nil
		} else {
			// Otherwise, delete the data in the cache and do another request
			intSettingCache.Delete(name)
		}
	}

	// Get the thing from the backend instead if not cached
	res, err := PostRequestBackendGeneric[getSettingIntResponse]("/node/get_int_setting", map[string]interface{}{
		"id":      node.ID,
		"token":   node.Token,
		"setting": name,
	})
	if err != nil {
		return 0, err
	}

	// Insert it into the cache
	boolSettingCache.Store(name, intSettingCacheData{
		Value:       res.Value,
		LastRequest: time.Now(),
	})

	// Parse the result
	return res.Value, err
}
