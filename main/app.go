package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Liphium/station/backend/database"
	backend_starter "github.com/Liphium/station/backend/starter"
	"github.com/Liphium/station/backend/util/auth"
	chatserver_starter "github.com/Liphium/station/chatserver/starter"
	"github.com/Liphium/station/main/integration"
	spacestation_starter "github.com/Liphium/station/spacestation/starter"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {

	// Check if there is an extra environment file given
	if len(os.Args) == 1 {
		// Load environment variables from default location
		if godotenv.Load() != nil {
			printWithPrefix("No .env file found")
		}
	} else {
		// Load environment variables from specified location
		if godotenv.Load(os.Args[1]) != nil {
			printWithPrefix("Specified file " + os.Args[1] + " not found")
		}
	}

	printWithPrefix("Starting Liphium station..")
	printWithPrefix("Starting backend..")

	// Check if a system uuid is set
	if os.Getenv("SYSTEM_UUID") == "" {
		printWithPrefix("Please set the following environment variables in your .env file:")
		printWithPrefix("SYSTEM_UUID=\"" + uuid.New().String() + "\"")

		return
	}

	// Start backend
	backend_starter.Startup(true)

	// Create default database stuff
	backend_starter.CreateDefaultObjects()
	if err := database.DBConn.Where("id = ?", 1).First(&database.App{}).Error; err != nil {
		if err := database.DBConn.Create(&database.App{
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

	if err := database.DBConn.Where("id = ?", 2).First(&database.App{}).Error; err != nil {
		if err := database.DBConn.Create(&database.App{
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
	var defaultChatNode database.Node
	if err := database.DBConn.Where("id = ? AND app_id = ?", 1, 1).First(&defaultChatNode).Error; err != nil {
		defaultChatNode = database.Node{
			ID:              1,
			AppID:           1,
			Load:            0,
			PeformanceLevel: 1,
			Token:           auth.GenerateToken(300),
			Domain:          os.Getenv("CHAT_NODE"),
			Status:          database.StatusStopped,
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

	var defaultSpaceNode database.Node
	if err := database.DBConn.Where("id = ? AND app_id = ?", 2, 2).First(&defaultSpaceNode).Error; err != nil {
		defaultSpaceNode = database.Node{
			ID:              2,
			AppID:           2,
			Load:            0,
			PeformanceLevel: 1,
			Token:           auth.GenerateToken(300),
			Domain:          os.Getenv("SPACE_NODE"),
			Status:          database.StatusStopped,
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
	worked := spacestation_starter.Start(false)
	if !worked {

		// Just block the main thread for infinity until paused (this should be enough)
		time.Sleep(time.Hour * 30000)
	}
}

func printWithPrefix(s string) {
	fmt.Println("[station] " + s)
}
