package nodes

import (
	"errors"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/node"
)

func Node(id uint, token string) (node.Node, error) {

	// Check if token is valid
	var found node.Node
	if err := database.DBConn.Where("id = ?", id).Take(&found).Error; err != nil {
		return node.Node{}, err
	}

	if found.Token != token {
		return node.Node{}, errors.New("invalid.token")
	}

	return found, nil
}
