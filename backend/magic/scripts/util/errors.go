package magic_util

import (
	"log"
	"runtime/debug"
	"testing"

	"gorm.io/gorm"
)

// Handles the error of the transaction :D
func DatabaseError(t *testing.T, tx *gorm.DB) {
	if tx.Error != nil {
		printError(t, "database error", tx.Error)
	}
}

func WebsocketError(t *testing.T, err error) {
	if err != nil {
		printError(t, "websocket error", err)
	}
}

func AccountServiceError(t *testing.T, err error) {
	if err != nil {
		printError(t, "account service error", err)
	}
}

// Under the hood function for printing the error (to handle everything consistently)
func printError(t *testing.T, kind string, err error) {
	if t != nil {
		debug.PrintStack()
		t.Fatalf("%s: %v", kind, err)
	} else {
		log.Fatalf("%s: %v", kind, err)
	}
}
