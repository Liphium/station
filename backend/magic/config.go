package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Liphium/magic/mconfig"
	backend_starter "github.com/Liphium/station/backend/starter"
)

// This is the function called once you run the project
func Run(ctx *mconfig.Context) {

	// Add the databases
	main := mconfig.NewPostgresDatabase("main")
	ctx.AddDatabase(main)

	// Allocate port for the backend
	basePort := ctx.ValuePort(3000)

	// Create the file store
	fileRepo := filepath.Join(ctx.MagicDirectory(), "files", ctx.Profile())
	if err := os.MkdirAll(fileRepo, 0755); err != nil {
		log.Fatalln("Couldn't create file directory for profile:", err)
	}

	// Clear the file repo in case in testing mode
	if ctx.Profile() == "test" {
		if err := os.RemoveAll(fileRepo); err != nil {
			log.Fatalln("Couldn't remove all files in test directory:", err)
		}
	}

	// Add the databases to the environment
	ctx.WithEnvironment(&mconfig.Environment{
		// Domain config (SHOULD NEVER CHANGE IN PRODUCTION)
		"BASE_PATH": mconfig.ValueWithBase(
			[]mconfig.EnvironmentValue{basePort},
			func(output []string) string {
				return fmt.Sprintf("localhost:%s", output[0])
			},
		),
		"BASE_PORT": basePort,

		// Backend configuration
		"APP_NAME":       mconfig.ValueStatic("Liphium"),
		"TESTING":        mconfig.ValueStatic("true"),
		"TESTING_AMOUNT": mconfig.ValueStatic("2"),
		"LISTEN":         mconfig.ValueStatic("127.0.0.1"),
		"PROTOCOL":       mconfig.ValueStatic("http://"),
		"SYSTEM_UUID":    mconfig.ValueStatic("fb2b217b-db14-4500-9b11-1dd675532e76"), // DO NOT USE THIS IN PRODUCTION
		"JWT_SECRET":     mconfig.ValueStatic("secret"),                               // DO NOT USE THIS IN PRODUCTION
		"SMTP_PRINT":     mconfig.ValueStatic("true"),                                 // DO NOT USE IN PRODUCTION

		// File storage location
		"FILE_REPO_TYPE": mconfig.ValueStatic("local"),
		"FILE_REPO":      mconfig.ValueStatic(fileRepo),

		// Database for backend
		"DB_USER":     main.Username(),
		"DB_PASSWORD": main.Password(),
		"DB_DATABASE": main.DatabaseName(ctx),
		"DB_HOST":     main.Host(ctx),
		"DB_PORT":     main.Port(ctx),
	})

	// Load secrets
	if err := ctx.LoadSecretsToEnvironment(".env"); err != nil {
		log.Fatalln(err)
	}
}

func Start() {
	backend_starter.Startup(false)
}
