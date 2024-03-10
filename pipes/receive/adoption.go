package receive

import (
	"errors"
	"log"
	"time"

	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/connection"
	"github.com/bytedance/sonic"
	"nhooyr.io/websocket"
)

func ReceiveWSAdoption(request string) (pipes.Node, error) {

	// Unmarshal
	var adoptionRq connection.AdoptionRequest
	err := sonic.Unmarshal([]byte(request), &adoptionRq)
	if err != nil {
		return pipes.Node{}, err
	}

	// Check token
	if adoptionRq.Token != pipes.CurrentNode.Token {
		return pipes.Node{}, errors.New("invalid token")
	}

	pipes.Log.Printf("[ws] Incoming event stream from node %s connected.", adoptionRq.Adopting.ID)
	pipes.AddNode(adoptionRq.Adopting)

	// Connect output stream (if not already connected)
	if !connection.ExistsWS(adoptionRq.Adopting.ID) {

		go func() {
			time.Sleep(2 * time.Second)

			log.Println("Connecting to", adoptionRq.Adopting.ID)

			connection.IterateWS(func(id string, value *websocket.Conn) bool {
				log.Println(id)
				return true
			})

			if err := connection.ConnectWS(adoptionRq.Adopting); err != nil {
				log.Println(err)
			}
		}()
	}

	return adoptionRq.Adopting, nil
}

func AdoptUDP(bytes []byte) error {

	// Remove adoption request prefix
	bytes = bytes[2:]

	// Unmarshal
	var adoptionRq connection.AdoptionRequest
	err := sonic.Unmarshal(bytes, &adoptionRq)
	if err != nil {
		return err
	}

	// Check token
	if adoptionRq.Token != pipes.CurrentNode.Token {
		return errors.New("invalid token")
	}

	pipes.Log.Printf("[udp] Incoming event stream from node %s connected.", adoptionRq.Adopting.ID)
	pipes.AddNode(adoptionRq.Adopting)

	// Connect output stream (if not already connected)
	if !connection.ExistsUDP(adoptionRq.Adopting.ID) {
		connection.ConnectUDP(adoptionRq.Adopting)
	}

	return nil
}
