package routes

import (
	"strings"
	"time"

	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	pipeshroutes "github.com/Liphium/station/pipeshandler/routes"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/studio"
	"github.com/Liphium/station/spacestation/util"
	"github.com/gofiber/fiber/v2"
)

// Node protocol version
const ProtocolVersion = 8

func SetupRoutes(router fiber.Router) {
	router.Post("/socketless", socketlessEvent)
	router.Post("/ping", ping)

	router.Post("/about", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"protocol_version": ProtocolVersion,
		})
	})

	// These are publicly accessible yk (cause this can be public information cause encryption and stuff)
	router.Post("/info", roomInfo)

	// Encrypted routes (at /enc to prevent issues)
	router.Route("/enc", encryptedRoutes)

	setupPipesFiber(router)
}

func encryptedRoutes(router fiber.Router) {

	// For joining a Space (no matter from where, used for decentralization and normal to make the API consistent)
	router.Post("/join", joinSpace)

	// For creating a Space and generating a connection token for it
	router.Post("/create", createSpace)
}

func setupPipesFiber(router fiber.Router) {
	caching.SSInstance = pipeshandler.Setup(pipeshandler.Config{
		Secret:              []byte(integration.JwtSecret),
		ExpectedConnections: 10_0_0_0,       // 10 thousand, but funny
		SessionDuration:     time.Hour * 24, // This is kinda important

		// Report nodes as offline
		NodeDisconnectHandler: func(node pipes.Node) {
			integration.ReportOffline(node)
		},

		// Handle client disconnect
		ClientDisconnectHandler: func(client *pipeshandler.Client) {
			if integration.Testing {
				util.Log.Println("Client disconnected:", client.ID)
			}

			// Remove from the table
			caching.LeaveTable(client.Session, client.ID)

			// Delete all the Warps the guy has
			caching.StopWarpsBy(client.Session, client.ID)

			// Disconnect the guy from studio
			studio.GetStudio(client.Session).Disconnect(client.ID)

			// Remove from room
			caching.RemoveMember(client.Session, client.ID)

			// Delete the studio in case the room has been deleted
			if _, ok := caching.GetRoom(client.Session); !ok {
				studio.DeleteStudio(client.Session)
			}
		},

		// Validate token and create room
		TokenValidateHandler: func(claims *pipeshandler.ConnectionTokenClaims, attachments string) bool {

			if !strings.HasPrefix(claims.Extra, "oj-") {
				util.Log.Println("no prefix for only join")
				return true
			}

			// Make sure the room exists
			claims.Session = strings.TrimPrefix(claims.Extra, "oj-") // Session is the room id and since that's now passed through extra we'll just set session to it
			_, valid := caching.GetRoom(claims.Session)
			if !valid {
				util.Log.Println("the room doesn't exist")
				return true
			}
			return false
		},

		// Handle enter network
		ClientConnectHandler: func(client *pipeshandler.Client, key string) bool {

			if integration.Testing {
				util.Log.Println("Client connected:", client.ID)
			}

			return false
		},

		// Handle client entering network
		ClientEnterNetworkHandler: func(client *pipeshandler.Client, key string) bool {
			return false
		},

		ErrorHandler: func(err error) {
			util.Log.Println("Pipes error:", err)
		},
	})
	router.Route("/", func(router fiber.Router) {
		pipeshroutes.SetupRoutes(router, caching.SSNode, caching.SSInstance, false)
	})
}
