package backend_starter

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Liphium/station/backend/database"
	routes_v1 "github.com/Liphium/station/backend/routes/v1"
	"github.com/Liphium/station/backend/util"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func Startup(routine bool) {

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		JSONEncoder:       sonic.Marshal,
		JSONDecoder:       sonic.Unmarshal,
		StreamRequestBody: true, // TODO: Proper request body protection (Make only certain endpoints accept streams)
	})

	util.TestAES()

	// Load environment variables (don't if isolated cause not needed)
	var err error
	if !routine {
		err = godotenv.Load()
		if err != nil {
			util.Log.Fatal("Error loading .env file")
		}
	}
	util.JWT_SECRET = os.Getenv("JWT_SECRET")

	// Set node protocol
	if os.Getenv("PROTOCOL") == "" {
		util.NodeProtocol = "https://"
	} else {
		util.NodeProtocol = os.Getenv("PROTOCOL")
	}

	// Connect to the databases
	database.Connect()

	app.Use(cors.New())
	app.Use(logger.New())

	// Handle routing
	app.Route("/", routes_v1.Router)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from the backend, you probably shouldn't be here though.. Anyway, enjoy your time!")
	})

	// Ask user for test mode
	testMode()

	// Listen on port 3000
	listenAddress := fmt.Sprintf("%s:%s", os.Getenv("LISTEN"), os.Getenv("BASE_PORT"))
	if os.Getenv("CLI") == "true" {
		go func() {
			err = app.Listen(listenAddress)

			util.Log.Println(err.Error())
		}()

		// Listen for commands
		if routine {
			go listenForCommands()
		} else {
			listenForCommands()
		}
	} else {

		var err error
		if routine {
			go app.Listen(listenAddress)
		} else {
			err = app.Listen(listenAddress)
		}

		if err != nil {
			panic(err)
		}
	}
}

func testMode() {

	if os.Getenv("TESTING") != "" {
		util.Testing = os.Getenv("TESTING") == "true"
		if util.Testing {
			util.Log.Println("Test mode enabled (Read from .env).")
		}
	} else {
		util.Log.Println("Do you want to continue in test mode? (y/n)")

		scanner := bufio.NewScanner(os.Stdin)

		scanner.Scan()
		util.Testing = scanner.Text() == "y"
	}

	if !util.Testing {
		return
	}

	util.Log.Println("Test mode enabled.")

	/* not need for now
	var foundNodes []node.Node
	database.DBConn.Find(&foundNodes)

	for _, n := range foundNodes {
		if n.Status == node.StatusStarted {
			util.Log.Println("Stopping node", n.Domain)

			nodes.TurnOff(&n, node.StatusStopped)
		}
	}
	*/
}
