package account

import (
	"time"
)

type Account struct {
	ID string `json:"id" gorm:"primaryKey"` // 8 character-long string

	Email     string    `json:"email" gorm:"uniqueIndex"`
	Username  string    `json:"username"`
	Tag       string    `json:"tag"`
	RankID    uint      `json:"rank"`
	CreatedAt time.Time `json:"created_at"`

	Rank           Rank             `json:"-" gorm:"foreignKey:RankID"`
	Authentication []Authentication `json:"-" gorm:"foreignKey:Account"`
	Sessions       []Session        `json:"-" gorm:"foreignKey:Account"`
}
