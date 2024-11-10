package verify

// Permissions
var Permissions = map[string]int16{
	"use_services": 10,
	"admin":        100,
}

// Permission names
const PermissionUseServices = "use_services"
const PermissionAdmin = "admin"

// Check if the current user has a specific permission
func (i *SessionInformation) HasPermission(perm string) bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	// Get the level of the specified permission
	permLevel, valid := Permissions[perm]
	if !valid {
		return false
	}

	return i.permissionLevel >= permLevel
}
