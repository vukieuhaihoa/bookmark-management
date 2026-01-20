package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system.
type User struct {
	ID          string    `gorm:"primaryKey;type:uuid;column:id" json:"id"`
	Username    string    `gorm:"unique;not null;column:username" json:"username"`
	Email       string    `gorm:"unique;not null;column:email" json:"email"`
	Password    string    `gorm:"not null;column:password" json:"-"`
	DisplayName string    `gorm:"column:display_name" json:"display_name"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName specifies the table name for the User model.
//
// Returns:
//   - string: The name of the database table for the User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook that is triggered before a new User record is created in the database.
// It generates a new UUID for the user ID if it is not already set.
//
// Parameters:
//   - tx: The GORM database transaction
//
// Returns:
//   - error: An error if UUID generation fails, otherwise nil
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
