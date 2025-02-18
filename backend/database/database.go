package database

import (
	"fmt"
	"os"
	"time"

	"github.com/Liphium/station/backend/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB

func Connect() {
	url := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_DATABASE") + " port=" + os.Getenv("DB_PORT")

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		util.Log.Fatal("Something went wrong during the connection with the database.", err)
	}

	util.Log.Println("Successfully connected to the database.")

	// Configure the database driver
	driver, _ := db.DB()

	driver.SetMaxIdleConns(10)
	driver.SetMaxOpenConns(100)
	driver.SetConnMaxLifetime(time.Hour)

	// Add the uuid extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		fmt.Println("uuid extension 'uuid-ossp' not found.")
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(
		// Account related tables
		&Account{},
		&Authentication{},
		&Session{},
		&Rank{},
		&PublicKey{},
		&ProfileKey{},
		&VaultKey{},
		&SignatureKey{},
		&StoredActionKey{},
		&CloudFile{},
		&Invite{},
		&InviteCount{},

		// Properties related tables
		&Friendship{},
		&Profile{},
		&StoredAction{},
		&AStoredAction{},
		&VaultEntry{},
		&KeyRequest{},

		// Node related tables
		&Node{},
		&NodeCreation{},

		// Server related tables
		&App{},
		&Setting{},
	)

	// Assign the database to the global variable
	DBConn = db
}
