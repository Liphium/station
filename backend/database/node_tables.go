package database

import (
	"time"
)

type NodeCreation struct {
	Token string    `json:"token" gorm:"primaryKey"`
	Date  time.Time `json:"date"`
}

type Node struct {
	ID uint `json:"id" gorm:"primaryKey"`

	AppID           uint    `json:"app"` // App ID
	Token           string  `json:"token"`
	Domain          string  `json:"domain"`
	Load            float64 `json:"load"`
	PeformanceLevel float32 `json:"performance_level"`

	// 1 - Started, 2 - Stopped, 3 - Error
	Status uint `json:"status"`
}

// Convert the node to a returnable entity
func (n *Node) ToEntity() NodeEntity {
	return NodeEntity{
		ID:     n.ID,
		Token:  n.Token,
		App:    n.AppID,
		Domain: n.Domain,
	}
}

const StatusStarted = 1
const StatusStopped = 2
const StatusError = 3

type NodeEntity struct {
	ID     uint   `json:"id"`
	Token  string `json:"token"`
	App    uint   `json:"app"`
	Domain string `json:"domain"`
}
