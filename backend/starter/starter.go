package backend_starter

import (
	"bufio"
	"log"
	"os"
	"time"

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

	log.SetOutput(os.Stdout)

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		JSONEncoder:       sonic.Marshal,
		JSONDecoder:       sonic.Unmarshal,
		StreamRequestBody: true, // TODO: Proper request body protection (Make only certain endpoints accept streams)
	})

	util.TestAES()

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
	app.Route("/v1", routes_v1.Router)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from the backend, you probably shouldn't be here though.. Anyway, enjoy your time!")
	})

	// Ask user for test mode
	testMode()

	// Listen on port 3000
	if os.Getenv("CLI") == "true" {
		go func() {
			err = app.Listen(os.Getenv("LISTEN"))

			log.Println(err.Error())
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
			go app.Listen(os.Getenv("LISTEN"))
		} else {
			err = app.Listen(os.Getenv("LISTEN"))
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
			log.Println("Test mode enabled (Read from .env).")
		}
	} else {
		log.Println("Do you want to continue in test mode? (y/n)")

		scanner := bufio.NewScanner(os.Stdin)

		scanner.Scan()
		util.Testing = scanner.Text() == "y"
	}

	if !util.Testing {
		return
	}

	log.Println("Test mode enabled.")

	token, _ := util.Token("123", "123", 100, time.Now().Add(time.Hour*24))

	log.Println("Test token: " + token)

	/* not need for now
	var foundNodes []node.Node
	database.DBConn.Find(&foundNodes)

	for _, n := range foundNodes {
		if n.Status == node.StatusStarted {
			log.Println("Stopping node", n.Domain)

			nodes.TurnOff(&n, node.StatusStopped)
		}
	}
	*/
}
