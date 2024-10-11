package database

// Apps created on the server
type App struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Tag         string `json:"tag"` // Application tag (for discovering if a certain app runs on an instance)
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     uint   `json:"version"`
	AccessLevel uint   `json:"access_level"`
}

// Settings in the admin panel
type Setting struct {
	Name  string `gorm:"primaryKey"`
	Value string
}
