package requests

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/google/uuid"
)

// GetSession gets the session from the database (returns false if it doesn't exist)
func GetSession(id uuid.UUID, session *account.Session) bool {

	if err := database.DBConn.Model(session).Where("id = ?", id).Take(&session).Error; err != nil {
		return false
	}

	return true
}

// CheckSessionPermission checks if the session has the required permission level (returns true if it doesn't)
func CheckSessionPermission(token string, permission uint, session *account.Session) bool {

	err := database.DBConn.Take(session, token).Error

	if err != nil || session.PermissionLevel < permission {
		return true
	}

	return false
}
