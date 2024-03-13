package spacestation_starter

import (
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/connection"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/caching/games/launcher"
	"github.com/Liphium/station/spacestation/handler"
	"github.com/Liphium/station/spacestation/routes"
	"github.com/Liphium/station/spacestation/server"
	"github.com/Liphium/station/spacestation/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Start() {

	// Setup memory
	caching.SetupMemory()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	if !integration.Setup(integration.IdentifierSpaceNode) {
		return
	}

	server.InitLiveKit()

	launcher.InitGames()
	nodeData := integration.Nodes[integration.IdentifierSpaceNode]
	pipes.SetupCurrent(fmt.Sprintf("%d", nodeData.NodeId), nodeData.NodeToken)
	util.Log.Println("Starting..")

	// Query current node AND JWT TOKEN
	_, _, currentApp, domain := integration.GetCurrent(integration.IdentifierChatNode)
	currentNodeData := integration.Nodes[integration.IdentifierChatNode]
	currentNodeData.AppId = currentApp
	integration.Nodes[integration.IdentifierChatNode] = currentNodeData

	// Setup routes (called here because of the jwt secret)
	app.Route("/", routes.SetupRoutes)

	util.Log.Printf("Node %s on app %d\n", pipes.CurrentNode.ID, currentApp)

	protocol := os.Getenv("WEBSOCKET_PROTOCOL")
	if protocol == "" {
		protocol = "wss://"
	}
	pipes.SetupWS(protocol + domain + "/connect")
	handler.Initialize()

	// Report online status
	res := integration.SetOnline(integration.IdentifierSpaceNode)
	parseNodes(res)

	// Check if test mode or production
	args := strings.Split(domain, ":")
	var err error
	util.Port, err = strconv.Atoi(os.Getenv("SPACE_NODE_PORT"))
	if err != nil {
		util.Log.Println("Error: Couldn't parse port of current node")
		return
	}
	util.UDPPort = util.Port + 1
	pipes.SetupUDP(fmt.Sprintf("%s:%d", args[0], util.UDPPort))

	// Test encryption
	first := testEncryption()
	second := testEncryption()

	if reflect.DeepEqual(first, second) {
		util.Log.Println("Error: Encryption is not working properly")
		return
	}

	util.Log.Println("Encryption is working properly!")

	pipes.DebugLogs = true

	// Create testing room
	if integration.Testing {
		caching.CreateRoom("id", "test")

		amount, err := strconv.Atoi(os.Getenv("TESTING_AMOUNT"))
		if err != nil {
			util.Log.Println("Error: Couldn't parse testing amount")
			return
		}

		for i := 0; i < amount; i++ {
			clientId := util.GenerateToken(5)
			connection := caching.EmptyConnection(clientId, "id")
			valid := caching.JoinRoom("id", connection.ClientID)
			if !valid {
				util.Log.Println("Error: Couldn't join room")
				return
			}
			util.Log.Println("--- TESTING CLIENT ---")
			util.Log.Println(connection.ClientID + ":" + connection.KeyBase64())
			util.Log.Println("----------------------")
		}
	}

	// Close caches on exit
	defer caching.CloseCaches()

	// Connect to other nodes
	pipes.IterateNodes(func(_ string, node pipes.Node) bool {

		util.Log.Println("Connecting to node " + node.WS)

		if err := connection.ConnectWS(node); err != nil {
			util.Log.Println(err.Error())
		}
		return true
	})

	// Start on localhost
	err = app.Listen(fmt.Sprintf("%s:%d", os.Getenv("LISTEN"), util.Port))
	if err != nil {
		panic(err)
	}
}

// This function is used to test if the encryption is working properly and always different
func testEncryption() []byte {

	encrypted, err := connection.Encrypt(pipes.CurrentNode.ID, []byte("H"))
	if err != nil {
		util.Log.Println("Error: Couldn't encrypt message")
		return nil
	}

	util.Log.Println("Encrypted message: " + base64.StdEncoding.EncodeToString(encrypted))

	decrypted, err := connection.Decrypt(pipes.CurrentNode.ID, encrypted)
	if err != nil {
		util.Log.Println("Error: Couldn't decrypt message")
		return nil
	}

	util.Log.Println("Decrypted message: " + string(decrypted))

	return encrypted
}

// Shared function between all nodes
func parseNodes(res map[string]interface{}) bool {

	if res["nodes"] == nil {
		return true
	}

	nodeList := res["nodes"].([]interface{})

	for _, node := range nodeList {
		n := node.(map[string]interface{})

		// Extract port and domain
		args := strings.Split(n["domain"].(string), ":")
		domain := args[0]
		port, err := strconv.Atoi(args[1])
		if err != nil {
			util.Log.Println("Error: Couldn't parse port of node " + n["id"].(string))
			return true
		}

		// Add node to pipes
		pipes.AddNode(pipes.Node{
			ID:    fmt.Sprintf("%d", int64(n["id"].(float64))),
			Token: n["token"].(string),
			WS:    "ws://" + fmt.Sprintf("%s:%d", domain, port) + "/adoption",
		})
	}

	return false
}
