package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Course struct {
	ID          string    `json:"id" gorm:"type:char(36);primaryKey"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
}

func (course *Course) BeforeCreate(tx *gorm.DB) (err error) {
	if course.ID == "" {
		course.ID = uuid.NewString()
	}

	return nil
}
