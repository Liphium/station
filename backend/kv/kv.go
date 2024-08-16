package kv

import "sync"

// The purpose of this package is to provide a wrapper for a potential key value
// store that could be used in the future when there are multiple instances of this
// application running. For now it's just a simple hash map which is gonna be fine
// for 99.9% of instances. Except when we want it to scale, which is why I'm preparing
// the backend for this possibility.

var cache *sync.Map = &sync.Map{}

// Get a value from the key value store
func Get(key string) (any, bool) {
	return cache.Load(key)
}

// Store a value in the key value store
func Store(key string, value any) {
	cache.Store(key, value)
}
