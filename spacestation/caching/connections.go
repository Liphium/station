package caching

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"net"

	"github.com/Liphium/station/spacestation/util"
	"github.com/dgraph-io/ristretto"
)

type Connection struct {
	ID             string
	Room           string
	CurrentSession string
	UDP            *net.UDPAddr
	Key            []byte
	Cipher         cipher.Block
}

func (c *Connection) KeyBase64() string {
	return base64.StdEncoding.EncodeToString(c.Key)
}

// ! Always use cost 1
var connectionsCache *ristretto.Cache // ConnectionID -> Connection

func setupConnectionsCache() {

	var err error
	connectionsCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 10_000_000, // 1 million expected connections
		MaxCost:     1 << 30,    // 1 GB
		BufferItems: 64,
	})

	if err != nil {
		panic(err)
	}
}

func EmptyConnection(connId string, room string) Connection {

	// Generate encryption key
	key, err := util.GenerateKey()
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Store in cache
	conn := Connection{
		ID:     connId,
		Room:   room,
		UDP:    nil,
		Key:    key,
		Cipher: block,
	}
	connectionsCache.Set(connId, conn, 1)

	return conn
}

func GetConnection(connId string) (Connection, bool) {
	conn, valid := connectionsCache.Get(connId)
	if !valid {
		return Connection{}, false
	}
	return conn.(Connection), valid
}

// TODO: Create a test for all deletion functions
func DeleteConnection(connId string) {
	connectionsCache.Del(connId)
}
