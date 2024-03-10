package app

type AppSetting struct {
	ID uint `json:"id" gorm:"primaryKey"`

	App   uint   `json:"app"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
