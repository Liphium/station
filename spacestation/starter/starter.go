package spacestation_starter

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/handler"
	"github.com/Liphium/station/spacestation/routes"
	"github.com/Liphium/station/spacestation/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Start(loadEnv bool) bool {

	// Setup memory
	caching.SetupMemory()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	if !integration.Setup(integration.IdentifierSpaceNode, loadEnv) {
		return false
	}

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
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from space station! This is the service that handles all the calling stuff (also tabletop). Since you're here you're probably trying some things! If you are, thank you, and please report security issues to Liphium if you find any. You can find us at https://liphium.com.")
	})

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
		return false
	}

	// Test encryption (to make sure they don't produce the same outcome)
	first := testEncryption()
	second := testEncryption()

	if reflect.DeepEqual(first, second) {
		util.Log.Println("Error: Encryption is not working properly")
		return false
	}

	util.Log.Println("Encryption is working properly!")

	pipes.DebugLogs = true

	// Create testing room
	if integration.Testing {
		caching.CreateRoom("id")

		amount, err := strconv.Atoi(os.Getenv("TESTING_AMOUNT"))
		if err != nil {
			util.Log.Println("Error: Couldn't parse testing amount")
			return false
		}

		for i := 0; i < amount; i++ {
			connID := util.GenerateToken(5)
			connection := caching.EmptyConnection(connID, "id")
			valid := caching.JoinRoom("id", connection.ID)
			if !valid {
				util.Log.Println("Error: Couldn't join room")
				return false
			}
			util.Log.Println("--- TESTING CLIENT ---")
			util.Log.Println(connection.ID + ":" + connection.KeyBase64())
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

	return true
}

// This function is used to test if the encryption is working properly and always different
func testEncryption() []byte {

	// Encrypt something
	encrypted, err := caching.SSNode.Encrypt(caching.SSNode.ID, []byte("Hello world"))
	if err != nil {
		util.Log.Println("Error: Couldn't encrypt message")
		return nil
	}

	// Test the decryption as well
	_, err = caching.SSNode.Decrypt(caching.SSNode.ID, encrypted)
	if err != nil {
		util.Log.Println("Error: Couldn't decrypt message")
		return nil
	}

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
