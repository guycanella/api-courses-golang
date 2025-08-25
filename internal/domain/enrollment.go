package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Enrollment struct {
	ID        string    `json:"id"         gorm:"type:char(36);primaryKey"`
	UserID    string    `json:"user_id"    gorm:"type:char(36);not null;index:uidx_user_course,unique"`
	CourseID  string    `json:"course_id"  gorm:"type:char(36);not null;index:uidx_user_course,unique"`
	CreatedAt time.Time `json:"created_at"`

	User   User   `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Course Course `json:"-" gorm:"foreignKey:CourseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (enrollment *Enrollment) BeforeCreate(tx *gorm.DB) (err error) {
	if enrollment.ID == "" {
		enrollment.ID = uuid.NewString()
	}

	return nil
}
