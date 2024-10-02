package nodes

import (
	"errors"

	"github.com/Liphium/station/backend/database"
)

func Node(id uint, token string) (database.Node, error) {

	// Check if token is valid
	var found database.Node
	if err := database.DBConn.Where("id = ?", id).Take(&found).Error; err != nil {
		return database.Node{}, err
	}

	if found.Token != token {
		return database.Node{}, errors.New("invalid.token")
	}

	return found, nil
}
