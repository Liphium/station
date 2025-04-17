package database

// This is only for storage and sync between devices
// TODO: Replace with a vault entry (or remove completely)
type Status struct {
	ID string `gorm:"primaryKey"` // Account ID

	Data string `gorm:"not null"` // Encrypted data
	Node uint   `gorm:"not null"`
}
