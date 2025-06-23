package chat_routes

import (
	"time"

	"github.com/Liphium/station/backend/chat"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/zapshare"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipeshandler"
	pipeshroutes "github.com/Liphium/station/pipeshandler/routes"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {

}

func Unauthorized(router fiber.Router) {

	// Create the gateway
	chat.Instance = pipeshandler.Setup(pipeshandler.Config{
		Secret:              []byte(util.JwtSecret),
		ExpectedConnections: 10_0_0_0,       // 10 thousand, but funny
		SessionDuration:     time.Hour * 24, // This is kinda important

		// Handle client disconnect
		ClientDisconnectHandler: func(client *pipeshandler.Client) {

			// Print debug stuff if in debug mode
			if integration.Testing {
				util.Log.Println("Client disconnected:", client.ID)
			}

			// Cancel all zap transactions
			zapshare.CancelTransactionByAccount(client.ID)

			// TODO: Maybe handle disconnections a little more?
		},

		// Handle token validation (nothing to do here)
		TokenValidateHandler: func(claims *pipeshandler.ConnectionTokenClaims, key string) bool {
			return false
		},

		// Handle enter network
		ClientConnectHandler: func(client *pipeshandler.Client, key string) bool {
			return false
		},

		// Handle client entering network
		ClientEnterNetworkHandler: func(client *pipeshandler.Client, key string) bool {
			if integration.Testing {
				util.Log.Println("Client connected:", client.ID)
			}
			return false
		},

		ErrorHandler: func(err error) {
			util.Log.Printf("pipeshandler error: %s \n", err.Error())
		},
	})

	// Add all the routes for the gateway
	router.Route("/", func(router fiber.Router) {
		pipeshroutes.SetupRoutes(router, nil, chat.Instance, false)
	})
}
