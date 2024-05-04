package account

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`

	Email     string    `json:"email" gorm:"uniqueIndex"`
	Username  string    `json:"username" gorm:"uniqueIndex"`
	RankID    uint      `json:"rank"`
	CreatedAt time.Time `json:"created_at"`

	Rank           Rank             `json:"-" gorm:"foreignKey:RankID"`
	Authentication []Authentication `json:"-" gorm:"foreignKey:Account"`
	Sessions       []Session        `json:"-" gorm:"foreignKey:Account"`
}
