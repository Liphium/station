package pipes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/cornelk/hashmap"
)

type Node struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	WS    string `json:"ws,omitempty"`  // Websocket ip
	UDP   string `json:"udp,omitempty"` // UDP ip
	SL    string `json:"sl,omitempty"`  // Socketless pipe

	// Encryption
	Cipher cipher.Block `json:"-"`
}

var Log = log.New(log.Writer(), "pipes ", log.Flags())

var nodes = hashmap.New[string, Node]()

var CurrentNode Node

func SetupCurrent(id string, token string) {

	if len(token) < 32 {
		panic("Token is too short (must be longer than 32 characters for AES-256 encryption)")
	}

	// Create encryption cipher
	tokenHash := sha256.Sum256([]byte(token))
	encryptionKey := tokenHash[:]

	Log.Println("Encryption key:", base64.StdEncoding.EncodeToString(encryptionKey))

	cipher, err := aes.NewCipher(encryptionKey)
	if err != nil {
		panic(err)
	}

	CurrentNode = Node{
		ID:     id,
		Token:  token,
		WS:     "",
		UDP:    "",
		SL:     "",
		Cipher: cipher,
	}
}

func SetupWS(ws string) {
	CurrentNode.WS = ws
}

func SetupUDP(udp string) {
	CurrentNode.UDP = udp
}

func SetupSocketless(sl string) {
	CurrentNode.SL = sl
}

func GetNode(id string) *Node {

	if id == CurrentNode.ID {
		return &CurrentNode
	}

	// Get node
	node, ok := nodes.Get(id)
	if !ok {
		return nil
	}

	return &node
}

func AddNode(node Node) {

	// Create encryption cipher
	tokenHash := sha256.Sum256([]byte(node.Token))
	encryptionKey := tokenHash[:]
	cipher, err := aes.NewCipher(encryptionKey)
	if err != nil {
		Log.Println("[node] Error adding node", node.ID, ":", err)
		return
	}

	node.Cipher = cipher
	nodes.Insert(node.ID, node)
}

func DeleteNode(node string) {
	nodes.Del(node)
}

// IterateConnections iterates over all connections. If the callback returns false, the iteration stops.
func IterateNodes(callback func(string, Node) bool) {
	nodes.Range(callback)
}
