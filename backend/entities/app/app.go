package app

type App struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Version     uint   `json:"version"`
	AccessLevel uint   `json:"access_level"`
}
