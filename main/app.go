package main

import (
	"fmt"
	"os"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/app"
	"github.com/Liphium/station/backend/entities/node"
	backend_starter "github.com/Liphium/station/backend/starter"
	"github.com/Liphium/station/backend/util/auth"
	chatserver_starter "github.com/Liphium/station/chatserver/starter"
	"github.com/Liphium/station/main/integration"
)

func main() {
	printWithPrefix("Starting Liphium station..")
	printWithPrefix("Starting backend..")

	pub, priv, err := backend_starter.GenerateKeyPair()
	if err != nil {
		printWithPrefix("Error generating key pair: " + err.Error())
		return
	}

	// Set environment variables
	os.Setenv("TC_PUBLIC_KEY", pub)
	os.Setenv("TC_PRIVATE_KEY", priv)

	// Start backend
	backend_starter.Startup(true)

	// Create default database stuff
	backend_starter.CreateDefaultObjects()
	if err := database.DBConn.Where("id = ?", 1).First(&app.App{}).Error; err != nil {
		if err := database.DBConn.Create(&app.App{
			ID:          1,
			Name:        "Chat",
			Description: "Chat application",
			Version:     0,
			AccessLevel: 20,
		}).Error; err != nil {
			panic(err)
		}
	}

	if err := database.DBConn.Where("id = ?", 2).First(&app.App{}).Error; err != nil {
		if err := database.DBConn.Create(&app.App{
			ID:          2,
			Name:        "Spaces",
			Description: "Spaces application",
			Version:     0,
			AccessLevel: 20,
		}).Error; err != nil {
			panic(err)
		}
	}

	os.Setenv("CHAT_APP", "1")
	os.Setenv("SPACES_APP", "2")

	// Create default nodes
	var defaultChatNode node.Node
	if err := database.DBConn.Where("id = ? AND app = ?", 1, 1).First(&defaultChatNode).Error; err != nil {
		if err := database.DBConn.Create(&node.Node{
			ID:              1,
			ClusterID:       1,
			AppID:           1,
			Load:            0,
			PeformanceLevel: 1,
			Token:           auth.GenerateToken(300),
			Domain:          os.Getenv("CHAT_NODE"),
			Status:          node.StatusStopped,
		}).Error; err != nil {
			panic(err)
		}
	}
	integration.Nodes[integration.IdentifierChatNode] = integration.NodeData{
		NodeToken: defaultChatNode.Token,
		NodeId:    defaultChatNode.ID,
		AppId:     defaultChatNode.AppID,
	}

	var defaultSpaceNode node.Node
	if err := database.DBConn.Where("id = ? AND app = ?", 2, 2).First(&defaultSpaceNode).Error; err != nil {
		if err := database.DBConn.Create(&node.Node{
			ID:              2,
			ClusterID:       1,
			AppID:           2,
			Load:            0,
			PeformanceLevel: 1,
			Token:           auth.GenerateToken(300),
			Domain:          os.Getenv("SPACE_NODE"),
			Status:          node.StatusStopped,
		}).Error; err != nil {
			panic(err)
		}
	}
	integration.Nodes[integration.IdentifierSpaceNode] = integration.NodeData{
		NodeToken: defaultSpaceNode.Token,
		NodeId:    defaultSpaceNode.ID,
		AppId:     defaultSpaceNode.AppID,
	}

	// Start chat server
	printWithPrefix("Starting chat server..")
	chatserver_starter.Start()
}

func printWithPrefix(s string) {
	fmt.Println("[station] " + s)
}
