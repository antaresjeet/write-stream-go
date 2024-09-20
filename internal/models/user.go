package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	SocialID  string    `gorm:"unique;not null"`
	Email     string    `gorm:"unique;not null"`
	Name      string    `gorm:"not null"`
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
