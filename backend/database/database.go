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

	// Migrate account related tables
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Authentication{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&Rank{})
	db.AutoMigrate(&PublicKey{})
	db.AutoMigrate(&ProfileKey{})
	db.AutoMigrate(&VaultKey{})
	db.AutoMigrate(&SignatureKey{})
	db.AutoMigrate(&StoredActionKey{})
	db.AutoMigrate(&CloudFile{})
	db.AutoMigrate(&Invite{})
	db.AutoMigrate(&InviteCount{})

	// Migrate account properties related tables
	db.AutoMigrate(&Friendship{})
	db.AutoMigrate(&Profile{})
	db.AutoMigrate(&StoredAction{})
	db.AutoMigrate(&AStoredAction{})
	db.AutoMigrate(&VaultEntry{})
	db.AutoMigrate(&KeyRequest{})

	// Migrate node related tables
	db.AutoMigrate(&Node{})
	db.AutoMigrate(&NodeCreation{})

	// Migrate server related tables
	db.AutoMigrate(&App{})
	db.AutoMigrate(&Setting{})

	// Assign the database to the global variable
	DBConn = db
}
