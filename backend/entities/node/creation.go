package node

import "time"

type NodeCreation struct {
	Token string    `json:"token" gorm:"primaryKey"`
	Date  time.Time `json:"date"`
}
