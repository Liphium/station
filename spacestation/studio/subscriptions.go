package studio

import (
	"sync"
)

type Subscription struct {
	mutex   *sync.Mutex
	client  string // Id of the client subscribing (read-only)
	track   string // Id of the track subscribed to (read-only)
	channel string // The channel the client is subscribed to
}

func (s *Subscription) Delete() {

}
