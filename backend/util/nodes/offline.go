package nodes

import (
	"log"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/node"
)

func TurnOff(node *node.Node, status uint) {

	node.Status = status
	node.Load = 0
	database.DBConn.Save(node)

	// Disconnect all sessions
	go DisconnectAll(node)
}

func DisconnectAll(node *node.Node) {

	// Disconnect all sessions
	database.DBConn.Model(&account.Session{}).Where("node = ?", node.ID).Updates(map[string]interface{}{
		"node": 0,
		"app":  0,
	})

	log.Println("Disconnected all sessions from node: " + node.Domain + "!")
}
