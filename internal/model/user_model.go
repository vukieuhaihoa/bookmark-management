package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          string    `gorm:"primaryKey;type:uuid;column:id" json:"id"`
	Username    string    `gorm:"unique;not null;column:username" json:"username"`
	Email       string    `gorm:"unique;not null;column:email" json:"email"`
	Password    string    `gorm:"not null;column:password" json:"-"`
	DisplayName string    `gorm:"column:display_name" json:"display_name"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
