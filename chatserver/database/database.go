package database

import (
	"os"
	"time"

	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/database/fetching"
	"github.com/Liphium/station/chatserver/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB

func Connect() {
	url := "host=" + os.Getenv("CN_DB_HOST") + " user=" + os.Getenv("CN_DB_USER") + " password=" + os.Getenv("CN_DB_PASSWORD") + " dbname=" + os.Getenv("CN_DB_DATABASE") + " port=" + os.Getenv("CN_DB_PORT")

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
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

	// Migrate the schema
	db.AutoMigrate(&conversations.Conversation{})
	db.AutoMigrate(&conversations.ConversationToken{})
	db.AutoMigrate(&conversations.Message{})
	db.AutoMigrate(&fetching.Status{})

	// Assign the database to the global variable
	DBConn = db
}
