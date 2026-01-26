package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/repository"
	"github.com/vukieuhaihoa/bookmark-management/pkg/jwtutils"
	"github.com/vukieuhaihoa/bookmark-management/pkg/stringutils"
)

const TokenExpirationDuration = 24 * time.Hour

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
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

	// Login authenticates a user with the provided username and password.
	// Returns a JWT token if authentication is successful, or an error if it fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - username: The username of the user attempting to log in.
	//   - password: The password of the user attempting to log in.
	//
	// Returns:
	//   - string: The JWT token if authentication is successful.
	//   - error: An error if authentication fails, otherwise nil.
	Login(ctx context.Context, username, password string) (string, error)

	// GetUserByID retrieves a user by their ID.
	// Returns the user or an error if the operation fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - id: The ID of the user to be retrieved.
	//
	// Returns:
	//   - *model.User: The user model if found.
	//   - error: An error if the retrieval fails or the user is not found.
	GetUserByID(ctx context.Context, id string) (*model.User, error)

	// UpdateUserByID updates a user's display name and email by their ID.
	// Returns an error if the operation fails.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - id: The ID of the user to be updated.
	//   - displayName: The new display name for the user.
	//   - email: The new email address for the user.
	//
	// Returns:
	//   - error: An error if the update fails, otherwise nil.
	UpdateUserByID(ctx context.Context, id, displayName, email string) error
}

type user struct {
	userRepo        repository.User
	passwordHashing stringutils.PasswordHashing
	jwtGenerator    jwtutils.JWTGenerator
}

// NewUser creates a new instance of the User service.
//
// Parameters:
//   - userRepo: The user repository used for database operations.
//   - passwordHashing: The password hashing utility for securing passwords.
//   - jwtGenerator: The JWT generator for creating authentication tokens.
//
// Returns:
//   - User: A new user service instance.
func NewUser(userRepo repository.User, passwordHashing stringutils.PasswordHashing, jwtGenerator jwtutils.JWTGenerator) User {
	return &user{
		userRepo:        userRepo,
		passwordHashing: passwordHashing,
		jwtGenerator:    jwtGenerator,
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
		return nil, err
	}

	return createdUser, nil
}

// Login authenticates a user with the provided username and password.
// If authentication is successful, it generates and returns a JWT token.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - username: The username of the user attempting to log in.
//   - password: The password of the user attempting to log in.
//
// Returns:
//   - string: The JWT token if authentication is successful.
//   - error: An error if authentication fails, otherwise nil.
func (u *user) Login(ctx context.Context, username, password string) (string, error) {
	user, err := u.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	ok := u.passwordHashing.CompareHashAndPassword(user.Password, password)
	if !ok {
		return "", ErrInvalidCredentials
	}

	jwtContent := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(TokenExpirationDuration).Unix(),
	}

	token, err := u.jwtGenerator.GenerateToken(jwtContent)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByID retrieves a user by their ID.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - id: The ID of the user to be retrieved.
//
// Returns:
//   - *model.User: The user model if found.
//   - error: An error if the retrieval fails or the user is not found.
func (u *user) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}

// UpdateUserByID updates a user's display name and email by their ID.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - id: The ID of the user to be updated.
//   - displayName: The new display name for the user.
//   - email: The new email address for the user.
//
// Returns:
//   - *model.User: The updated user model.
//   - error: An error if the update fails, otherwise nil.
func (u *user) UpdateUserByID(ctx context.Context, id, displayName, email string) error {
	updatedUser := &model.User{
		DisplayName: displayName,
		Email:       email,
	}

	return u.userRepo.UpdateUserByID(ctx, id, updatedUser)
}
