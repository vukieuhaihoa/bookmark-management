package repository

import (
	"context"
	"fmt"

	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
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

	// GetUserByUsername retrieves a user from the database by their username.
	// Returns the user or an error if the operation fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - username: The username of the user to be retrieved.
	//
	// Returns:
	//   - *model.User: The user model if found.
	//   - error: An error if the retrieval fails or the user is not found.
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

	// GetUserByID retrieves a user from the database by their ID.
	// Returns the user or an error if the operation fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - id: The ID of the user to be retrieved.
	//
	// Returns:
	//   - *model.User: The user model if found.
	//   - error: An error if the retrieval fails or the user is not found.
	GetUserByID(ctx context.Context, id string) (*model.User, error)

	// UpdateUserByID updates an existing user in the database by their ID.
	// Returns an error if the operation fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - id: The ID of the user to be updated.
	//   - updatedUser: The user model containing the updated details.
	//
	// Returns:
	//   - error: An error if the update fails, otherwise nil.
	UpdateUserByID(ctx context.Context, id string, updatedUser *model.User) error
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
		return nil, dbutils.CatchDBError(err)
	}

	return newUser, nil
}

// GetUserByUsername retrieves a user from the database by their username.
// It takes a context and a username as input and returns the user or an error.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - username: The username of the user to be retrieved.
//
// Returns:
//   - *model.User: The user model if found.
//   - error: An error if the retrieval fails or the user is not found.
func (u *user) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return u.GetUserByField(ctx, "username", username)
}

// GetUserByID retrieves a user from the database by their ID.
// It takes a context and an ID as input and returns the user or an error.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - id: The ID of the user to be retrieved.
//
// Returns:
//   - *model.User: The user model if found.
//   - error: An error if the retrieval fails or the user is not found.
func (u *user) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return u.GetUserByField(ctx, "id", id)
}

// GetUserByField retrieves a user from the database by a specified field and value.
// It takes a context, a field name, and a value as input and returns the user or an error.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - field: The field name to search by (e.g., "email", "username").
//   - value: The value of the field to match.
//
// Returns:
//   - *model.User: The user model if found.
//   - error: An error if the retrieval fails or the user is not found.
func (u *user) GetUserByField(ctx context.Context, field string, value string) (*model.User, error) {
	user := &model.User{}
	err := u.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(user).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return user, nil
}

func (u *user) UpdateUserByID(ctx context.Context, id string, updatedUser *model.User) error {
	result := u.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updatedUser)
	if result.Error != nil {
		return dbutils.CatchDBError(result.Error)
	}

	if result.RowsAffected == 0 {
		return dbutils.ErrRecordNotFoundType
	}

	return nil
}
