package node

type Cluster struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Name    string `json:"name" gorm:"unique"`
	Country string `json:"country" gorm:"unique"`
}
