package database

import (
	"fmt"
	"os"
	"time"

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

	// Add the uuid extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		fmt.Println("uuid extension 'uuid-ossp' not found.")
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&Conversation{})
	db.AutoMigrate(&ConversationToken{})
	db.AutoMigrate(&Message{})
	db.AutoMigrate(&Status{})

	// Assign the database to the global variable
	DBConn = db
}
