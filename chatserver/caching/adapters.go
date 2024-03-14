package caching

import (
	"errors"

	"github.com/Liphium/station/chatserver/util"
	"github.com/dgraph-io/ristretto"
)

// ! Always use cost 1
var adapterCache *ristretto.Cache // Account ID -> All adapters linked (to remove them later)

func setupAdapterCache() {
	var err error

	// TODO: Check if values really are enough
	adapterCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of objects expected (1,000,000).
		MaxCost:     1 << 30, // maximum cost of cache (1,000,000).
		BufferItems: 64,      // something from the docs
	})
	if err != nil {
		panic(err)
	}
}

// Insert all the adapters
func InsertAdapters(account string, adapters []string) {
	adapterCache.Set(account, adapters, 1)
}

// Delete all adapters for an account
func DeleteAdapters(account string) error {

	// Get all adapters
	obj, valid := adapterCache.Get(account)
	if !valid {
		return errors.New("not found")
	}
	adapters := obj.([]string)

	// Remove adapters from pipes
	for _, adapterName := range adapters {
		util.Log.Println("DELETED " + adapterName)
		Node.RemoveAdapterWS(adapterName)
	}

	adapterCache.Del(account)
	return nil
}
