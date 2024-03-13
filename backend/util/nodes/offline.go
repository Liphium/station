package nodes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/entities/node"
	"github.com/Liphium/station/backend/util"
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

	util.Log.Println("Disconnected all sessions from node: " + node.Domain + "!")
}
