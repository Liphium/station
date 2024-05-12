package routes

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/fetching"
	conversation_routes "github.com/Liphium/station/chatserver/routes/conversations"
	"github.com/Liphium/station/chatserver/routes/ping"
	zapshare_routes "github.com/Liphium/station/chatserver/routes/zapshare"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	pipesfroutes "github.com/Liphium/station/pipeshandler/routes"
	"github.com/bytedance/sonic"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {

	// Return the public key for TC Protection
	router.Post("/pub", func(c *fiber.Ctx) error {

		// Return the public key in a packaged form (string)
		return c.JSON(fiber.Map{
			"pub": integration.PackageRSAPublicKey(integration.NodePublicKey),
		})
	})

	router.Post("/ping", ping.Pong)

	// Pipes fiber doesn't need(/support) encrypted routes (it actually does for socketless, which is why we now have a seperate )
	setupPipesFiber(router)

	// We don't need to encrypt the liveshare routes
	router.Route("/liveshare", zapshare_routes.Unencrypted)

	router.Route("/auth", authorizedRoutes)
	router.Route("/", encryptedRoutes)
}

func authorizedRoutes(router fiber.Router) {
	authorize(router)
	router.Route("/liveshare", zapshare_routes.Authorized)
}

func encryptedRoutes(router fiber.Router) {

	// Add Through Cloudflare Protection middleware
	router.Use(func(c *fiber.Ctx) error {

		// Get the AES encryption key from the Auth-Tag header
		aesKeyEncoded, valid := c.GetReqHeaders()["Auth-Tag"]
		if !valid {
			util.Log.Println("no header")
			return c.SendStatus(fiber.StatusPreconditionFailed)
		}
		aesKeyEncrypted, err := base64.StdEncoding.DecodeString(aesKeyEncoded[0])
		if err != nil {
			util.Log.Println("no decoding")
			return c.SendStatus(fiber.StatusPreconditionFailed)
		}

		// Decrypt the AES key using the private key of this node
		aesKey, err := integration.DecryptRSA(integration.NodePrivateKey, aesKeyEncrypted)
		if err != nil {
			return c.SendStatus(fiber.StatusPreconditionRequired)
		}

		// Decrypt the request body using the key attached to the Auth-Tag header
		decrypted, err := integration.DecryptAES(aesKey, c.Body())
		if err != nil {
			return c.SendStatus(fiber.StatusNetworkAuthenticationRequired)
		}

		// Set some variables for use when sending back the response
		c.Locals("body", decrypted)
		c.Locals("key", aesKey)
		c.Locals("srv_pub", integration.NodePrivateKey)

		// Go to the next middleware/handler
		return c.Next()
	})

	// No authorization needed for this route
	router.Post("/adoption/socketless", socketless)

	// Authorized by using a remote id or normal token
	authorize(router)

	// Authorized routes (for accounts with remote id only)
	router.Route("/conversations", conversation_routes.SetupRoutes)
}

func authorize(router fiber.Router) {
	// Authorized by using a remote id or normal token
	router.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS512,
			Key:    []byte(integration.JwtSecret),
		},

		// Checks if the token is expired
		SuccessHandler: func(c *fiber.Ctx) error {

			// Check if the JWT is expired
			if util.IsExpired(c) {
				return integration.InvalidRequest(c, "expired jwt token")
			}

			// Go to the next middleware/handler
			return c.Next()
		},

		// Error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			util.Log.Println(c.Route().Path, "jwt error:", err.Error())

			// Return error message
			return c.SendStatus(fiber.StatusUnauthorized)
		},
	}))
}

func setupPipesFiber(router fiber.Router) {
	caching.CSInstance = pipeshandler.Setup(pipeshandler.Config{
		Secret:              []byte(integration.JwtSecret),
		ExpectedConnections: 10_0_0_0,       // 10 thousand, but funny
		SessionDuration:     time.Hour * 24, // This is kinda important

		// Report nodes as offline
		NodeDisconnectHandler: func(node pipes.Node) {

			// Report that a node is offline to the backend
			integration.ReportOffline(node)
		},

		// Handle client disconnect
		ClientDisconnectHandler: func(client *pipeshandler.Client) {

			// Print debug stuff if in debug mode
			if integration.Testing {
				util.Log.Println("Client disconnected:", client.ID)
			}

			// Remove all adapters from pipes
			err := caching.DeleteAdapters(client.ID)
			if err != nil {
				util.Log.Println("COULDN'T DELETE ADAPTERS:", err.Error())
			}

			// Tell the backend that someone disconnected
			nodeData := integration.Nodes[integration.IdentifierChatNode]
			integration.PostRequest("/node/disconnect", map[string]interface{}{
				"node":    nodeData.NodeId,
				"token":   nodeData.NodeToken,
				"session": client.Session,
			})
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

			// Just for debug purposes
			util.Log.Println(base64.StdEncoding.EncodeToString(aesKey))

			// Set AES key in client data
			client.Data = ExtraClientData{aesKey}
			caching.CSInstance.UpdateClient(client)

			// Initialize the user and check if he needs to be disconnected
			disconnect := !initializeUser(client)
			util.Log.Println("Setup finish")
			if disconnect {
				util.Log.Println("Something went wrong with setup: ", client.ID)
			}
			return disconnect
		},

		//* Set the decoding middleware to use encryption
		ClientEncodingMiddleware: EncryptionClientEncodingMiddleware,
		DecodingMiddleware:       EncryptionDecodingMiddleware,

		ErrorHandler: func(err error) {
			util.Log.Printf("pipes-fiber error: %s \n", err.Error())
		},
	})
	router.Route("/", func(router fiber.Router) {
		pipesfroutes.SetupRoutes(router, caching.CSNode, caching.CSInstance, false)
	})
}

// Extra client data attached to the pipes-fiber client
type ExtraClientData struct {
	Key []byte // AES encryption key
}

// Middleware for pipes-fiber to add encryption support
func EncryptionDecodingMiddleware(client *pipeshandler.Client, instance *pipeshandler.Instance, bytes []byte) (pipeshandler.Message, error) {

	// Handle potential errors
	defer func() {
		if err := recover(); err != nil {
			instance.ReportClientError(client, "encryption failure", errors.ErrUnsupported)
		}
	}()

	// Decrypt the message using AES
	key := client.Data.(ExtraClientData).Key
	messageEncoded, err := integration.DecryptAES(key, bytes)
	if err != nil {
		return pipeshandler.Message{}, err
	}

	// Unmarshal the message using sonic
	var message pipeshandler.Message
	err = sonic.Unmarshal(messageEncoded, &message)
	if err != nil {
		return pipeshandler.Message{}, err
	}

	return message, nil
}

// Middleware for pipes-fiber to add encryption support (in encoding)
func EncryptionClientEncodingMiddleware(client *pipeshandler.Client, instance *pipeshandler.Instance, message []byte) ([]byte, error) {

	// Handle potential errors (with casting in particular)
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
	util.Log.Println("ENCODING KEY: "+base64.StdEncoding.EncodeToString(key), client.ID, string(message))
	result, err := integration.EncryptAES(key, message)
	hash := sha256.Sum256(result)
	util.Log.Println("hash: " + base64.StdEncoding.EncodeToString(hash[:]))
	return result, err
}

func initializeUser(client *pipeshandler.Client) bool {
	account := client.ID

	// Check if the account is already in the database
	var status fetching.Status
	if database.DBConn.Where(&fetching.Status{ID: account}).Take(&status).Error != nil {

		// Create a new status
		if database.DBConn.Create(&fetching.Status{
			ID:   account,
			Data: "-", // Status is disabled
			Node: integration.Nodes[integration.IdentifierChatNode].NodeId,
		}).Error != nil {
			return false
		}
	} else {

		// Update the status
		database.DBConn.Model(&fetching.Status{}).Where("id = ?", account).Update("node", util.NodeTo64(caching.CSNode.ID))
	}

	// Send current status
	caching.CSInstance.SendEvent(client, pipes.Event{
		Name: "setup_st",
		Data: map[string]interface{}{
			"data": status.Data,
			"node": status.Node,
		},
	})

	// Send the setup complete event
	caching.CSInstance.SendEvent(client, pipes.Event{
		Name: "setup_fin",
		Data: map[string]interface{}{},
	})

	return true
}
