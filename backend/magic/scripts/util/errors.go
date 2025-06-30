package magic_util

import (
	"log"

	"gorm.io/gorm"
)

// Handles the error of the transaction :D
func DatabaseError(tx *gorm.DB) {
	if tx.Error != nil {
		log.Fatalln("database error:", tx.Error)
	}
}

func WebsocketError(err error) {
	if err != nil {
		log.Fatalln("websocket error:", err)
	}
}
