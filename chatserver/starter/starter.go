package chatserver_starter

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/handler"
	"github.com/Liphium/station/chatserver/routes"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Start(routine bool) {

	fmt.Println("IF YOU ARE ON LINUX, MAKE SURE TO RUN THIS PROGRAM WITH RIGHT PERMISSIONS TO NODE_ENV")
	util.Log.SetOutput(os.Stdout)

	// Setting up the node
	if !integration.Setup(integration.IdentifierChatNode, !routine) {
		return
	}

	// Setup environment
	allowUnsafe := os.Getenv("CN_ALLOW_UNSAFE")
	if allowUnsafe == "" {
		util.AllowUnsafe = false
	} else if allowUnsafe == "true" {
		util.AllowUnsafe = true
	}
	chatNodePath := os.Getenv("CHAT_NODE")
	if chatNodePath == "" || strings.Contains(chatNodePath, "http://") || strings.Contains(chatNodePath, "https://") {
		panic("Please set the CHAT_NODE environment variable to the domain of the chat server. WITHOUT http:// or https://, specify that in the PROTOCOL environment variable.")
	}
	util.OwnPath = os.Getenv("PROTOCOL") + os.Getenv("CHAT_NODE")

	// Connect to the database
	database.Connect()
	caching.SetupCaches()

	// Create fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	// Query current node
	_, _, currentApp, domain := integration.GetCurrent(integration.IdentifierChatNode)
	currentNodeData := integration.Nodes[integration.IdentifierChatNode]
	currentNodeData.AppId = currentApp
	integration.Nodes[integration.IdentifierChatNode] = currentNodeData

	nodeData := integration.Nodes[integration.IdentifierChatNode]
	caching.CSNode = pipes.SetupCurrent(fmt.Sprintf("%d", nodeData.NodeId), nodeData.NodeToken)

	// Report online status
	res := integration.SetOnline(integration.IdentifierChatNode)
	parseNodes(res)

	util.Log.Printf("Node %s on app %d\n", caching.CSNode.ID, currentApp)

	caching.CSNode.SetupSocketless(domain + "/adoption/socketless")

	app.Use(logger.New(logger.Config{
		Format: "chat | " + logger.ConfigDefault.Format,
	}))
	app.Route("/", routes.Setup)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from the chat server, this is handling all your messages. Since you're here you're probably trying some things! If you are, thank you, and please report security issues to Liphium if you find any. You can find us at https://liphium.com.")
	})

	// Create handlers
	handler.Create()

	// Check if test mode or production
	var port int
	var err error
	port, err = strconv.Atoi(os.Getenv("CHAT_NODE_PORT"))
	if err != nil {
		panic(err)
	}

	protocol := os.Getenv("WEBSOCKET_PROTOCOL")
	if protocol == "" {
		protocol = "wss://"
	}
	caching.CSNode.SetupWS(protocol + domain + "/connect")

	// Connect to other nodes
	caching.CSNode.IterateNodes(func(_ string, node pipes.Node) bool {

		util.Log.Println("Connecting to node " + node.WS)

		if err := caching.CSNode.ConnectToNodeWS(node); err != nil {
			util.Log.Println(err.Error())
		}
		return true
	})

	pipes.DebugLogs = true

	// Start on localhost
	if routine {
		go app.Listen(fmt.Sprintf("%s:%d", os.Getenv("LISTEN"), port))
	} else {
		app.Listen(fmt.Sprintf("%s:%d", os.Getenv("LISTEN"), port))
	}
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
		caching.CSNode.AddNode(pipes.Node{
			ID:    fmt.Sprintf("%d", int64(n["id"].(float64))),
			Token: n["token"].(string),
			WS:    "ws://" + fmt.Sprintf("%s:%d", domain, port) + "/adoption",
		})
	}

	return false
}
