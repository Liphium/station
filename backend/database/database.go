package database

import (
	"fmt"
	"os"
	"time"

	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/entities/app"
	"github.com/Liphium/station/backend/entities/node"
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

	// Migrate account related tables
	db.AutoMigrate(&account.Account{})
	db.AutoMigrate(&account.Authentication{})
	db.AutoMigrate(&account.Session{})
	db.AutoMigrate(&account.Rank{})
	db.AutoMigrate(&account.PublicKey{})
	db.AutoMigrate(&account.ProfileKey{})
	db.AutoMigrate(&account.SignatureKey{})
	db.AutoMigrate(&account.StoredActionKey{})
	db.AutoMigrate(&account.CloudFile{})
	db.AutoMigrate(&account.Invite{})
	db.AutoMigrate(&account.InviteCount{})

	// Migrate account properties related tables
	db.AutoMigrate(&properties.Friendship{})
	db.AutoMigrate(&properties.Profile{})
	db.AutoMigrate(&properties.StoredAction{})
	db.AutoMigrate(&properties.AStoredAction{})
	db.AutoMigrate(&properties.VaultEntry{})
	db.AutoMigrate(&properties.KeyRequest{})

	// Migrate node related tables
	db.AutoMigrate(&node.Cluster{})
	db.AutoMigrate(&node.Node{})
	db.AutoMigrate(&node.NodeCreation{})

	// Migrate app related tables
	db.AutoMigrate(&app.App{})
	db.AutoMigrate(&app.AppSetting{})

	// Assign the database to the global variable
	DBConn = db
}
