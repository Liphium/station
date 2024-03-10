package service

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/fetching"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/integration"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

func User(client *pipeshandler.Client) bool {
	account := client.ID

	// Check if the account is already in the database
	var status fetching.Status
	if database.DBConn.Where(&fetching.Status{ID: account}).Take(&status).Error != nil {

		// Create a new status
		if database.DBConn.Create(&fetching.Status{
			ID:   account,
			Data: "-", // Status is disabled
			Node: integration.NODE_ID,
		}).Error != nil {
			return false
		}
	} else {

		// Update the status
		database.DBConn.Model(&fetching.Status{}).Where("id = ?", account).Update("node", util.NodeTo64(pipes.CurrentNode.ID))
	}

	// Send current status
	client.SendEvent(pipes.Event{
		Name: "setup_st",
		Data: map[string]interface{}{
			"data": status.Data,
			"node": status.Node,
		},
	})

	// Send the setup complete event
	client.SendEvent(pipes.Event{
		Name: "setup_fin",
		Data: map[string]interface{}{},
	})

	return true
}
