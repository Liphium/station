package routes

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	pipeshroutes "github.com/Liphium/station/pipeshandler/routes"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/handler"
	"github.com/Liphium/station/spacestation/util"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	router.Post("/socketless", socketlessEvent)
	router.Post("/ping", ping)

	router.Post("/pub", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"pub": integration.PackageRSAPublicKey(integration.NodePublicKey),
		})
	})

	// These are publicly accessible yk (cause this can be public information cause encryption and stuff)
	router.Post("/leave", leaveRoom)
	router.Post("/info", roomInfo)

	setupPipesFiber(router)
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

			// Remove from room
			caching.RemoveMember(client.Session, client.ID)
			caching.DeleteConnection(client.ID)

			// Send leave event
			handler.SendRoomData(client.Session)
		},

		// Validate token and create room
		TokenValidateHandler: func(claims *pipeshandler.ConnectionTokenClaims, attachments string) bool {

			// Create room (if needed)
			_, valid := caching.GetRoom(claims.Session)
			if !valid {
				util.Log.Println("Creating new room for", claims.Account, "("+claims.Session+")")
				caching.CreateRoom(claims.Session)
			} else {
				util.Log.Println("Room already exists for", claims.Account, "("+claims.Session+")")
			}

			return false
		},

		// Handle enter network
		ClientConnectHandler: func(client *pipeshandler.Client, key string) bool {

			// Get the AES key from attachments
			aesKeyEncrypted, err := base64.StdEncoding.DecodeString(key)
			if err != nil {
				return true
			}

			// Decrypt AES key
			aesKey, err := integration.DecryptRSA(integration.NodePrivateKey, aesKeyEncrypted)
			if err != nil {
				return true
			}

			// Set AES key in client data
			client.Data = ExtraClientData{aesKey}
			caching.SSInstance.UpdateClient(client)

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

		// Set default encoding middleware
		DecodingMiddleware:       EncryptionDecodingMiddleware,
		ClientEncodingMiddleware: EncryptionClientEncodingMiddleware,
	})
	router.Route("/", func(router fiber.Router) {
		pipeshroutes.SetupRoutes(router, caching.SSNode, caching.SSInstance, false)
	})
}

// Extra client data attached to the pipes-fiber client
type ExtraClientData struct {
	Key []byte // AES encryption key
}

// Middleware for pipes-fiber to add encryption support
func EncryptionDecodingMiddleware(client *pipeshandler.Client, instance *pipeshandler.Instance, bytes []byte) ([]byte, error) {

	// Handle potential errors
	defer func() {
		if err := recover(); err != nil {
			instance.ReportClientError(client, "encryption failure", errors.ErrUnsupported)
		}
	}()

	// Decrypt the message using AES
	key := client.Data.(ExtraClientData).Key
	return integration.DecryptAES(key, bytes)
}

// Middleware for pipes-fiber to add encryption support (in encoding)
func EncryptionClientEncodingMiddleware(client *pipeshandler.Client, instance *pipeshandler.Instance, message []byte) ([]byte, error) {

	// Handle potential errors
	defer func() {
		if err := recover(); err != nil {
			instance.ReportClientError(client, "encryption failure", errors.ErrUnsupported)
		}
	}()

	// Check if the encryption key is set
	if client.Data == nil {
		return nil, errors.New("no encryption key set")
	}

	// Encrypt the message using the client encryption key
	key := client.Data.(ExtraClientData).Key
	result, err := integration.EncryptAES(key, message)
	return result, err
}
