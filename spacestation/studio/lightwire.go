package studio

import (
	"log"
	"sync"

	"github.com/pion/webrtc/v4"
)

type Lightwire struct {
	client  *Client
	mutex   *sync.Mutex
	channel *webrtc.DataChannel
}

// Initialize the lightwire server implementation
func (lw *Lightwire) Init() {
	lw.mutex.Lock()
	defer lw.mutex.Unlock()

	// Set a higher buffer threshold (default is typically 0)
	lw.channel.SetBufferedAmountLowThreshold(256 * 1024) // 256KB

	// Add a listener to monitor buffer status
	lw.channel.OnBufferedAmountLow(func() {
		logger.Println(lw.client.id, "lightwire buffer amount low")
	})

	// Send all packets from lightwire to everyone in Studio
	lw.channel.OnMessage(func(msg webrtc.DataChannelMessage) {
		// Close the connection in case a string was sent
		if msg.IsString {
			log.Println("string message, closing")
			lw.Close()
			return
		}

		// Add the client id to the packet
		// Format: | id_length (8 bytes) | client_id (length of id_length) | voice_data (rest) |
		clientIdAsBytes := []byte(lw.client.id)
		newPacket := append([]uint8{uint8(len(clientIdAsBytes))}, clientIdAsBytes...)
		newPacket = append(newPacket, msg.Data...)

		// Forward the packet to all lightwire clients
		lw.client.studio.ForwardLightwirePacket(lw.client.id, newPacket)
	})
}

// Forward a packet to the lightwire data channel.
//
// Closes the connection in case of a failure.
func (lw *Lightwire) SendPacket(packet []byte) {
	lw.mutex.Lock()

	// Send the packet
	if err := lw.channel.Send(packet); err != nil {
		logger.Println(lw.client.id, "lightwire sending failure:", err)

		// Close the connection in case of failure
		lw.mutex.Unlock()
		lw.Close()
		return
	}

	lw.mutex.Unlock()
}

// Close the lightwire connection
func (lw *Lightwire) Close() {
	lw.mutex.Lock()
	defer lw.mutex.Unlock()
	lw.channel.Close()
	lw.client.handleLightwireClose()
}
