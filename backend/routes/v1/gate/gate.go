package chat_routes

import (
	"github.com/Liphium/station/backend/service"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/zapshare"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/neogate"
	"github.com/gofiber/fiber/v2"
)

func Authorized(router fiber.Router) {

}

func Unauthorized(router fiber.Router) {

	// Create the gateway
	service.Instance = neogate.Setup(neogate.Config{
		Secret: []byte(util.JwtSecret),

		// Handle client disconnect
		ClientDisconnectHandler: func(client *neogate.Client) {

			// Print debug stuff if in debug mode
			if integration.Testing {
				util.Log.Println("Client disconnected:", client.ID)
			}

			// Cancel all zap transactions
			zapshare.CancelTransactionByAccount(client.ID)

			// TODO: Maybe handle disconnections a little more?
		},

		// Handle enter network
		ClientConnectHandler: func(client *neogate.Client, key string) bool {
			return false
		},

		// Handle client entering network
		ClientEnterNetworkHandler: func(client *neogate.Client, key string) bool {
			if integration.Testing {
				util.Log.Println("Client connected:", client.ID)
			}
			return false
		},

		// Set the adapter name of the client to include the address
		ClientAdapterHandler: func(client *neogate.Client) string {
			return standards.LiphiumAddress(client.ID)
		},

		ErrorHandler: func(err error) {
			util.Log.Printf("pipeshandler error: %s \n", err.Error())
		},
	})

	// Add all the routes for the gateway
	router.Route("/connect", service.Instance.MountGateway)
}
