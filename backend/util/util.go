package util

import (
	"log"
	"os"

	"github.com/google/uuid"
)

// Environment variables
const EnvAppName = "APP_NAME" // Configure the app name

// Important variables
const ProtocolVersion = 8

var Testing = false
var LogErrors = true

var Log = log.New(os.Stdout, "backend ", log.Flags())

var JwtSecret = ""

// Get the system uuid set in the environment variables
func GetSystemUUID() uuid.UUID {
	id, err := uuid.Parse(os.Getenv("SYSTEM_UUID"))
	if err != nil {
		panic("Please set the SYSTEM_UUID env property.")
	}

	return id
}
