package caching

import "sync"

// Conversation id -> sync.Map of SharedSpace instances
var sharedSpacesMap = &sync.Map{}

type SharedSpace struct {
	Id      string
	Mutex   *sync.Mutex
	Members []string // Encrypted (member ids)
	Data    string   // Encrypted (name and stuff)
}
