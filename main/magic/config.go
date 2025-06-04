package config

import (
	"fmt"
	"log"

	"github.com/Liphium/magic/mconfig"
)

// This is the function called once you run the project
func Run(ctx *mconfig.Context) {

	// Add the databases
	main := mconfig.NewPostgresDatabase("main")
	chat := mconfig.NewPostgresDatabase("chat")
	ctx.AddDatabase(main)
	ctx.AddDatabase(chat)

	// Add the databases to the environment
	ctx.WithEnvironment(&mconfig.Environment{
		// Database for backend
		"DB_USER":     main.Username(),
		"DB_PASSWORD": main.Password(),
		"DB_DATABASE": main.DatabaseName(ctx),
		"DB_HOST":     main.Host(ctx),
		"DB_PORT":     main.Port(ctx),

		// Database for chatserver
		"CN_DB_USER":     main.Username(),
		"CN_DB_PASSWORD": main.Password(),
		"CN_DB_DATABASE": main.DatabaseName(ctx),
		"CN_DB_HOST":     main.Host(ctx),
		"CN_DB_PORT":     main.Port(ctx),
	})

	// Load secrets
	log.Println("profile:", ctx.Profile())
	if err := ctx.LoadSecretsToEnvironment(".env"); err != nil {
		log.Fatalln(err)
	}
}

func Start() {
	fmt.Println("Hello magic!")
}
