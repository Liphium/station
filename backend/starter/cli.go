package backend_starter

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/entities/app"
	"github.com/Liphium/station/backend/entities/node"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"gorm.io/gorm"
)

func listenForCommands() {
	for {
		fmt.Print("node-backend > ")
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "exit":
			os.Exit(0)
		case "create-default":
			CreateDefaultObjects()

		case "create-app":

			fmt.Print("App name: ")
			appName, _ := reader.ReadString('\n')
			appName = strings.TrimSpace(appName)

			fmt.Print("App description: ")
			appDescription, _ := reader.ReadString('\n')
			appDescription = strings.TrimSpace(appDescription)

			fmt.Print("App access level: ")
			appAccessLevel, _ := reader.ReadString('\n')
			appAccessLevel = strings.TrimSpace(appAccessLevel)

			// Create app
			accessLevel, err := strconv.Atoi(appAccessLevel)
			if err != nil {
				fmt.Println("Invalid access level")
				continue
			}

			app := &app.App{
				Name:        appName,
				Description: appDescription,
				Version:     0,
				AccessLevel: uint(accessLevel),
			}
			database.DBConn.Create(&app)

			fmt.Println("Created app with ID", app.ID)

		case "increment-version":

			fmt.Print("App id: ")
			appIdString, _ := reader.ReadString('\n')
			appIdString = strings.TrimSpace(appIdString)
			appId, err := strconv.Atoi(appIdString)
			if err != nil {
				fmt.Println("Please enter an integer as the ID.")
				continue
			}

			var application app.App
			if database.DBConn.Where("id = ?", appId).Take(&application).Error != nil {
				fmt.Println("Couldn't get the app from the database")
				continue
			}

			// Increment version
			if database.DBConn.Model(&app.App{}).Where("id = ?", application.ID).Update("version", application.Version+1).Error != nil {
				fmt.Println("Couldn't increment version of the app in the database")
				continue
			}

			fmt.Println("Version has been incremented!")

		case "create-node":

			// Generate new node token
			tk := auth.GenerateToken(100)

			// Save
			if err := database.DBConn.Create(&node.NodeCreation{
				Token: tk,
				Date:  time.Now(),
			}).Error; err != nil {
				fmt.Println("Failed to create node token")
				continue
			}

			fmt.Println("Created node token", tk)

		case "delete-data":

			fmt.Print("Account E-Mail: ")
			email, _ := reader.ReadString('\n')
			email = strings.TrimSpace(email)

			// Delete all data
			var acc account.Account
			if err := database.DBConn.Where("email = ?", email).Take(&acc).Error; err != nil {
				fmt.Println("Failed to find account")
				continue
			}

			database.DBConn.Where("account = ?", acc.ID).Delete(&account.Session{})
			database.DBConn.Where("id = ?", acc.ID).Delete(&account.ProfileKey{})
			database.DBConn.Where("id = ?", acc.ID).Delete(&account.StoredActionKey{})
			database.DBConn.Where("id = ?", acc.ID).Delete(&account.PublicKey{})
			database.DBConn.Where("id = ?", acc.ID).Delete(&account.SignatureKey{})
			database.DBConn.Where("account = ?", acc.ID).Delete(&properties.AStoredAction{})
			database.DBConn.Where("account = ?", acc.ID).Delete(&properties.StoredAction{})
			database.DBConn.Where("account = ?", acc.ID).Delete(&properties.Friendship{})
			database.DBConn.Where("account = ?", acc.ID).Delete(&properties.VaultEntry{})
			database.DBConn.Where("id = ?", acc.ID).Delete(&properties.Profile{})

		case "account-token":

			fmt.Print("Account E-Mail: ")
			email, _ := reader.ReadString('\n')
			email = strings.TrimSpace(email)

			// Get account
			var acc account.Account
			if err := database.DBConn.Where("email = ?", email).Preload("Rank").Take(&acc).Error; err != nil {
				fmt.Println("Failed to find account")
				continue
			}

			// Generate new token
			token, err := util.Token("test", acc.ID, acc.Rank.Level, time.Now().Add(time.Hour*24*365))
			if err != nil {
				fmt.Println("Failed to generate token")
				continue
			}

			fmt.Println("Token:", token)

		case "keypair":

			priv, pub, err := util.GenerateRSAKey(util.StandardKeySize)
			if err != nil {
				fmt.Println("Failed to generate a keypair!")
				continue
			}

			fmt.Println("Packaged public key:", util.PackageRSAPublicKey(pub))
			fmt.Println("Packaged private key:", util.PackageRSAPrivateKey(priv))

		case "test-message":

			fmt.Print("Test message: ")
			msg, _ := reader.ReadString('\n')
			msg = strings.TrimSpace(msg)

			_, pub, err := util.GenerateRSAKey(util.StandardKeySize)
			if err != nil {
				fmt.Println("Failed to generate a keypair!")
				continue
			}

			// Get default private and public key
			serverPub, err := util.UnpackageRSAPublicKey(os.Getenv("TC_PUBLIC_KEY"))
			if err != nil {
				panic("Couldn't unpackage public key. Required for v1 API. Please set TC_PUBLIC_KEY in your environment variables or .env file.")
			}

			encrypted, err := util.EncryptRSA(serverPub, []byte(msg))
			if err != nil {
				fmt.Println("Couldn't encrypt using server pub")
				continue
			}

			directory := os.Getenv("TC_WRITE_TO")
			if directory == "" {
				fmt.Println("Please provide the environment variable TC_WRITE_TO (file directory for message).")
				continue
			}
			directory = directory + "/test.msg"
			util.Log.Println(directory)
			err = os.WriteFile(directory, encrypted, os.ModeAppend)
			if err != nil {
				fmt.Println("Couldn't write file:", err)
				continue
			}

			fmt.Println("Packaged public key:", util.PackageRSAPublicKey(pub))

		case "invite-wave":

			invites := 100
			for invites > 0 {

				// Get a random account
				var acc account.Account
				if err := database.DBConn.Order("random()").Take(&acc).Error; err != nil {
					fmt.Println("Couldn't get a random account:", err.Error())
					break
				}

				// Get the current invite count
				var inviteCount account.InviteCount
				err := database.DBConn.Where("account = ?", acc.ID).Take(&inviteCount).Error
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					fmt.Println("Couldn't get invite count:", err.Error())
					break
				}

				if errors.Is(err, gorm.ErrRecordNotFound) {
					inviteCount.Count = 0
					inviteCount.Account = acc.ID
				}

				// Give the user 3 invites
				invites -= 3
				inviteCount.Count += 3
				if err := database.DBConn.Save(&inviteCount).Error; err != nil {
					fmt.Println("Couldn't increment invite count by 3:", err.Error())
					break
				}

				fmt.Println("Gave 3 invites to", acc.Email, "("+acc.ID.String()+")")
			}

			fmt.Println("Invite wave finished. Hope everyone enjoys them!")
			continue

		case "test-account":

			fmt.Print("Name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			acc := &account.Account{
				Email:    name + "@liphium.app",
				Username: name,
				RankID:   1, // Default
			}
			if err := database.DBConn.Create(&acc).Error; err != nil {
				fmt.Println("error:", err.Error())
				continue
			}

			hash, err := auth.HashPassword("yourmum123", acc.ID)
			if err != nil {
				return
			}

			if err := database.DBConn.Create(&account.Authentication{
				ID:      auth.GenerateToken(5),
				Account: acc.ID,
				Type:    account.TypePassword,
				Secret:  hash,
			}).Error; err != nil {
				fmt.Println("error:", err.Error())
				continue
			}

			fmt.Println("Name:", name)
			fmt.Println("Email:", name+"@liphium.app")
			fmt.Println("Password:", "yourmum123")

		case "generate-invite":

			invite := account.Invite{
				ID:      auth.GenerateToken(32),
				Creator: util.GetSystemUUID(),
			}
			if err := database.DBConn.Create(&invite).Error; err != nil {
				fmt.Println("err:", err.Error())
			}

			fmt.Println("invite:", invite.ID)

		case "help":
			fmt.Println("exit - Exit the application")
			fmt.Println("create-default - Create default ranks and cluster")
			fmt.Println("create-app - Create a new app")
			fmt.Println("increment-version - Increment the version of an app (when a breaking change is made)")
			fmt.Println("create-node - Get a node token (rest of setup in the CLI of the node)")
			fmt.Println("delete-data - Delete the data to restart the setup process on an account")
			fmt.Println("account-token - Generate a JWT token for an account")
			fmt.Println("keypair - Generate a new RSA key pair")
			fmt.Println("test-message - Encrypt a test message to send to an endpoint using TC")
			fmt.Println("invite-wave - Give out 100 random invites.")
			fmt.Println("generate-invite - Generate an invite.")
			fmt.Println("test-account - Create a test account.")

		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}

func GenerateKeyPair() (publicKey string, privateKey string, theError error) {
	priv, pub, err := util.GenerateRSAKey(util.StandardKeySize)
	if err != nil {
		return "", "", err
	}

	return util.PackageRSAPublicKey(pub), util.PackageRSAPrivateKey(priv), nil
}

func CreateDefaultObjects() {
	if database.DBConn.Where("name = ?", "Default").Take(&account.Rank{}).RowsAffected > 0 {
		fmt.Println("Default stuff already exists")
		return
	}

	// Create default ranks
	database.DBConn.Create(&account.Rank{
		Name:  "Default",
		Level: 20,
	})
	database.DBConn.Create(&account.Rank{
		Name:  "Admin",
		Level: 100,
	})

	// Create default cluster
	database.DBConn.Create(&node.Cluster{
		Name:    "Default cluster",
		Country: "DE",
	})

	fmt.Println("Created default ranks and cluster")
}
