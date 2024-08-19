package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/app"
	"github.com/Liphium/station/backend/entities/node"
	backend_starter "github.com/Liphium/station/backend/starter"
	"github.com/Liphium/station/backend/util/auth"
	chatserver_starter "github.com/Liphium/station/chatserver/starter"
	"github.com/Liphium/station/main/integration"
	spacestation_starter "github.com/Liphium/station/spacestation/starter"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {

	if godotenv.Load() != nil {
		printWithPrefix("No .env file found")
		return
	}

	printWithPrefix("Starting Liphium station..")
	printWithPrefix("Starting backend..")

	// Set environment variables
	if os.Getenv("TC_PUBLIC_KEY") == "" {
		pub, priv, err := backend_starter.GenerateKeyPair()
		if err != nil {
			printWithPrefix("Error generating key pair: " + err.Error())
			return
		}

		printWithPrefix("Please set the following environment variables in your .env file:")
		printWithPrefix("TC_PUBLIC_KEY=" + pub)
		printWithPrefix("TC_PRIVATE_KEY=" + priv)

		return
	}

	// Check if a system uuid is set
	if os.Getenv("SYSTEM_UUID") == "" {
		printWithPrefix("Please set the following environment variables in your .env file:")
		printWithPrefix("SYSTEM_UUID=" + uuid.New().String())

		return
	}

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
			Tag:         "liphium_chat",
		}).Error; err != nil {
			panic(err)
		}
		printWithPrefix("Created default chat app")
	}

	if err := database.DBConn.Where("id = ?", 2).First(&app.App{}).Error; err != nil {
		if err := database.DBConn.Create(&app.App{
			ID:          2,
			Name:        "Spaces",
			Description: "Spaces application",
			Version:     0,
			AccessLevel: 20,
			Tag:         "liphium_spaces",
		}).Error; err != nil {
			panic(err)
		}
		printWithPrefix("Created default spaces app")
	}

	// Start default nodes (spaces and chat)
	time.Sleep(time.Millisecond * 500)
	printWithPrefix("Starting default nodes..")

	os.Setenv("CHAT_APP", "1")
	os.Setenv("SPACES_APP", "2")

	// Create default nodes
	var defaultChatNode node.Node
	if err := database.DBConn.Where("id = ? AND app_id = ?", 1, 1).First(&defaultChatNode).Error; err != nil {
		defaultChatNode = node.Node{
			ID:              1,
			ClusterID:       1,
			AppID:           1,
			Load:            0,
			PeformanceLevel: 1,
			Token:           auth.GenerateToken(300),
			Domain:          os.Getenv("CHAT_NODE"),
			Status:          node.StatusStopped,
		}
		if err := database.DBConn.Create(&defaultChatNode).Error; err != nil {
			panic(err)
		}
		printWithPrefix("Created default chat node")
	}
	integration.Nodes[integration.IdentifierChatNode] = integration.NodeData{
		NodeToken: defaultChatNode.Token,
		NodeId:    defaultChatNode.ID,
		AppId:     defaultChatNode.AppID,
	}

	var defaultSpaceNode node.Node
	if err := database.DBConn.Where("id = ? AND app_id = ?", 2, 2).First(&defaultSpaceNode).Error; err != nil {
		defaultSpaceNode = node.Node{
			ID:              2,
			ClusterID:       1,
			AppID:           2,
			Load:            0,
			PeformanceLevel: 1,
			Token:           auth.GenerateToken(300),
			Domain:          os.Getenv("SPACE_NODE"),
			Status:          node.StatusStopped,
		}
		if err := database.DBConn.Create(&defaultSpaceNode).Error; err != nil {
			panic(err)
		}
		printWithPrefix("Created default space node")
	}
	integration.Nodes[integration.IdentifierSpaceNode] = integration.NodeData{
		NodeToken: defaultSpaceNode.Token,
		NodeId:    defaultSpaceNode.ID,
		AppId:     defaultSpaceNode.AppID,
	}

	// Start chat server
	printWithPrefix("Starting chat server..")
	chatserver_starter.Start(true)

	// Start space station
	printWithPrefix("Starting space station..")
	spacestation_starter.Start()
}

func printWithPrefix(s string) {
	fmt.Println("[station] " + s)
}
