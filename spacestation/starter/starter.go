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
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/caching/games/launcher"
	"github.com/Liphium/station/spacestation/handler"
	"github.com/Liphium/station/spacestation/routes"
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

	caching.InitLiveKit()

	launcher.InitGames()
	util.Log.Println("Starting..")

	// Query current node AND JWT TOKEN
	_, _, currentApp, domain := integration.GetCurrent(integration.IdentifierSpaceNode)
	currentNodeData := integration.Nodes[integration.IdentifierSpaceNode]
	currentNodeData.AppId = currentApp
	integration.Nodes[integration.IdentifierSpaceNode] = currentNodeData

	nodeData := integration.Nodes[integration.IdentifierSpaceNode]
	caching.SSNode = pipes.SetupCurrent(fmt.Sprintf("%d", nodeData.NodeId), nodeData.NodeToken)
	util.Log.Println("NODE", caching.SSNode.ID)

	// Setup routes (called here because of the jwt secret)
	app.Route("/", routes.SetupRoutes)

	util.Log.Printf("Node %s on app %d\n", caching.SSNode.ID, currentApp)

	protocol := os.Getenv("WEBSOCKET_PROTOCOL")
	if protocol == "" {
		protocol = "wss://"
	}
	caching.SSNode.SetupWS(protocol + domain + "/connect")
	handler.Initialize()

	// Report online status
	res := integration.SetOnline(integration.IdentifierSpaceNode)
	parseNodes(res)

	// Check if test mode or production
	var err error
	util.Port, err = strconv.Atoi(os.Getenv("SPACE_NODE_PORT"))
	if err != nil {
		util.Log.Println("Error: Couldn't parse port of current node")
		return
	}

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
		caching.CreateRoom("id")

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
	caching.SSNode.IterateNodes(func(_ string, node pipes.Node) bool {

		util.Log.Println("Connecting to node " + node.WS)

		if err := caching.SSNode.ConnectToNodeWS(node); err != nil {
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

	encrypted, err := caching.SSNode.Encrypt(caching.SSNode.ID, []byte("H"))
	if err != nil {
		util.Log.Println("Error: Couldn't encrypt message")
		return nil
	}

	util.Log.Println("Encrypted message: " + base64.StdEncoding.EncodeToString(encrypted))

	decrypted, err := caching.SSNode.Decrypt(caching.SSNode.ID, encrypted)
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
		caching.SSNode.AddNode(pipes.Node{
			ID:    fmt.Sprintf("%d", int64(n["id"].(float64))),
			Token: n["token"].(string),
			WS:    "ws://" + fmt.Sprintf("%s:%d", domain, port) + "/adoption",
		})
	}

	return false
}
