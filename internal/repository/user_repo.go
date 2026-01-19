package repository

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	"gorm.io/gorm"
)

// User represents the interface for user repository operations.
//
//go:generate mockery --name=User --filename=user_repo.go --output=./mocks
type User interface {
	// CreateUser creates a new user in the database.
	// Returns the created user or an error if the operation fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - user: The user model containing the details of the user to be created.
	//
	// Returns:
	//   - *model.User: The created user model.
	//   - error: An error if the creation fails, otherwise nil.
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
}

type user struct {
	db *gorm.DB
}

// NewUser creates a new instance of the User repository.
//
// Parameters:
//   - db: The GORM database connection.
//
// Returns:
//   - User: A new user repository instance.
func NewUser(db *gorm.DB) User {
	return &user{db: db}
}

// CreateUser creates a new user in the database.
// It takes a context and a user model as input and returns the created user or an error.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - newUser: The user model containing the details of the user to be created.
//
// Returns:
//   - *model.User: The created user model.
//   - error: An error if the creation fails, otherwise nil.
func (u *user) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	err := u.db.WithContext(ctx).Create(newUser).Error
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
