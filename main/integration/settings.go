package integration

import (
	"sync"
	"time"

	"github.com/Liphium/station/pipes"
)

// Predefined settings
const SettingDecentralizationEnabled = "decentralization.enabled"
const SettingDecentralizationUnsafeAllowed = "decentralization.allow_unsafe"

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
