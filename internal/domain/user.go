package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `json:"id" gorm:"type:char(36);primaryKey"`
	Email     string    `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	return nil
}

func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	return nil
}
