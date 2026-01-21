package service

import (
	"context"
	"errors"

	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/repository"
	"github.com/vukieuhaihoa/bookmark-management/pkg/stringutils"
)

var (
	ErrCannotCreateUser = errors.New("cannot create user")
)

// User represents the interface for user service operations.
//
//go:generate mockery --name=User --filename=user_service.go --output=./mocks
type User interface {
	// CreateUser creates a new user with the provided information.
	// Returns the created user or an error if the operation fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - username: The username of the new user.
	//   - password: The password of the new user.
	//   - displayName: The display name of the new user.
	//   - email: The email address of the new user.
	//
	// Returns:
	//   - *model.User: The created user model.
	//   - error: An error if the creation fails, otherwise nil.
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)
}

type user struct {
	userRepo        repository.User
	passwordHashing stringutils.PasswordHashing
}

// NewUser creates a new instance of the User service.
//
// Parameters:
//   - userRepo: The user repository used for database operations.
//   - passwordHashing: The password hashing utility for securing passwords.
//
// Returns:
//   - User: A new user service instance.
func NewUser(userRepo repository.User, passwordHashing stringutils.PasswordHashing) User {
	return &user{
		userRepo:        userRepo,
		passwordHashing: passwordHashing,
	}
}

// CreateUser creates a new user with the provided information.
// It hashes the password before storing the user in the database.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - username: The username of the new user.
//   - password: The password of the new user.
//   - displayName: The display name of the new user.
//   - email: The email address of the new user.
//
// Returns:
//   - *model.User: The created user model.
//   - error: An error if the creation fails, otherwise nil.
func (u *user) CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error) {
	hashedPassword, err := u.passwordHashing.Hash(password)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		Username:    username,
		Password:    hashedPassword,
		DisplayName: displayName,
		Email:       email,
	}

	createdUser, err := u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, ErrCannotCreateUser
	}

	return createdUser, nil
}
