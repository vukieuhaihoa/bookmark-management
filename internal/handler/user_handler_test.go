package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/service"
	svcMocks "github.com/vukieuhaihoa/bookmark-management/internal/service/mocks"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
	"github.com/vukieuhaihoa/bookmark-management/pkg/validators"
)

func init() {
	// Register custom validators before running tests
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password_strength", validators.PasswordStrength)
	}
}

func TestUser_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputRequest *createUserRequest

		setupRequest func(ctx *gin.Context, inputRequest *createUserRequest)

		setupMockSvc func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.User

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful create user",

			inputRequest: &createUserRequest{
				Username:    "testuser",
				Password:    "my_SECURE_password123@",
				DisplayName: "Test User",
				Email:       "testuser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createUserRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("CreateUser", ctx, inputRequest.Username, inputRequest.Password, inputRequest.DisplayName, inputRequest.Email).
					Return(&model.User{
						Username:    inputRequest.Username,
						DisplayName: inputRequest.DisplayName,
						Email:       inputRequest.Email,
						ID:          "de305d54-75b4-431b-adb2-eb6b9e546099",
						CreatedAt:   fixture.TestTime,
						UpdatedAt:   fixture.TestTime,
					}, nil)
				return mockUserSvc
			},

			expectedCode:     http.StatusCreated,
			expectedResponse: `{"data":{"id":"de305d54-75b4-431b-adb2-eb6b9e546099","username":"testuser","email":"testuser@example.com","display_name":"Test User","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"},"message":"Register an user successfully!"}`,
		},
		{
			name: "invalid request body",

			inputRequest: &createUserRequest{
				Username:    "",
				Password:    "short",
				DisplayName: "Test User",
				Email:       "invalid-email-format",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createUserRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.User {
				return svcMocks.NewUser(t) // No expectations since service should not be called
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid input fields","details":["Username is invalid (required)","Password is invalid (min)","Email is invalid (email)"]}`,
		},
		{
			name: "invalid request body - weak password",

			inputRequest: &createUserRequest{
				Username:    "testuser",
				Password:    "shortshort",
				DisplayName: "Test User",
				Email:       "testuser@gmail.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createUserRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.User {
				return svcMocks.NewUser(t) // No expectations since service should not be called
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid input fields","details":["Password is invalid (password_strength)"]}`,
		},
		{
			name: "duplicate username or email",

			inputRequest: &createUserRequest{
				Username:    "existinguser",
				Password:    "my_SECURE_password123@",
				DisplayName: "Existing User",
				Email:       "existinguser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createUserRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("CreateUser", mock.Anything, inputRequest.Username, inputRequest.Password, inputRequest.DisplayName, inputRequest.Email).
					Return(nil, dbutils.ErrDuplicationType)
				return mockUserSvc
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"username or email already exists"}`,
		},
		{
			name: "service layer error",

			inputRequest: &createUserRequest{
				Username:    "testuser",
				Password:    "my_SECURE_password123@",
				DisplayName: "Test User",
				Email:       "testuser@gmail.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *createUserRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("CreateUser", mock.Anything, inputRequest.Username, inputRequest.Password, inputRequest.DisplayName, inputRequest.Email).
					Return(nil, assert.AnError)
				return mockUserSvc
			},

			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"message":"Internal server error"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx, tc.inputRequest)
			mockUserSvc := tc.setupMockSvc(ctx, tc.inputRequest)

			userHandler := NewUser(mockUserSvc)
			userHandler.CreateUser(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestUser_Login(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputRequest *loginRequest
		setupRequest func(ctx *gin.Context, inputRequest *loginRequest)

		setupMockSvc func(ctx *gin.Context, inputRequest *loginRequest) *svcMocks.User

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful login",
			inputRequest: &loginRequest{
				Username: "testuser",
				Password: "my_SECURE_password123@",
			},
			setupRequest: func(ctx *gin.Context, inputRequest *loginRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(ctx *gin.Context, inputRequest *loginRequest) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("Login", mock.Anything, inputRequest.Username, inputRequest.Password).
					Return("mocked-jwt-token", nil)
				return mockUserSvc
			},
			expectedCode:     http.StatusOK,
			expectedResponse: `{"data":"mocked-jwt-token","message":"Logged in successfully!"}`,
		},
		{
			name: "invalid request body",
			inputRequest: &loginRequest{
				Username: "",
				Password: "",
			},
			setupRequest: func(ctx *gin.Context, inputRequest *loginRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(ctx *gin.Context, inputRequest *loginRequest) *svcMocks.User {
				return svcMocks.NewUser(t) // No expectations since service should not be called
			},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid input fields","details":["Username is invalid (required)","Password is invalid (required)"]}`,
		},
		{
			name: "password too short",
			inputRequest: &loginRequest{
				Username: "testuser",
				Password: "short",
			},
			setupRequest: func(ctx *gin.Context, inputRequest *loginRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(ctx *gin.Context, inputRequest *loginRequest) *svcMocks.User {
				return svcMocks.NewUser(t) // No expectations since service should not be called
			},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid input fields","details":["Password is invalid (gte)"]}`,
		},
		{
			name: "invalid credentials",
			inputRequest: &loginRequest{
				Username: "testuser",
				Password: "wrong_password",
			},
			setupRequest: func(ctx *gin.Context, inputRequest *loginRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(ctx *gin.Context, inputRequest *loginRequest) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("Login", mock.Anything, inputRequest.Username, inputRequest.Password).
					Return("", service.ErrInvalidCredentials)
				return mockUserSvc
			},
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"invalid username or password"}`,
		},
		{
			name: "service layer error",
			inputRequest: &loginRequest{
				Username: "testuser",
				Password: "my_SECURE_password123@",
			},
			setupRequest: func(ctx *gin.Context, inputRequest *loginRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(ctx *gin.Context, inputRequest *loginRequest) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("Login", mock.Anything, inputRequest.Username, inputRequest.Password).
					Return("", assert.AnError)
				return mockUserSvc
			},
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"message":"Internal server error"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx, tc.inputRequest)
			mockUserSvc := tc.setupMockSvc(ctx, tc.inputRequest)

			userHandler := NewUser(mockUserSvc)
			userHandler.Login(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestUser_GetProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)

		setupMockSvc func(ctx *gin.Context) *svcMocks.User

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful get profile",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user by setting userID in context
				ctx.Set("userID", "de305d54-75b4-431b-adb2-eb6b9e546099")
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("GetUserByID", mock.Anything, "de305d54-75b4-431b-adb2-eb6b9e546099").
					Return(&model.User{
						ID:          "de305d54-75b4-431b-adb2-eb6b9e546099",
						Username:    "testuser",
						Email:       "testuser@example.com",
						DisplayName: "Test User",
						CreatedAt:   fixture.TestTime,
						UpdatedAt:   fixture.TestTime,
					}, nil)
				return mockUserSvc
			},

			expectedCode:     http.StatusOK,
			expectedResponse: `{"data":{"id":"de305d54-75b4-431b-adb2-eb6b9e546099","username":"testuser","email":"testuser@example.com","display_name":"Test User","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"},"message":"User profile retrieved successfully!"}`,
		},
		{
			name: "unauthenticated request",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// No userID set in context to simulate unauthenticated request
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.User {
				return svcMocks.NewUser(t) // No expectations since service should not be called
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "invalid userID in context",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Set invalid userID in context
				ctx.Set("userID", 12345) // should be a string
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.User {
				return svcMocks.NewUser(t) // No expectations since service should not be called
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "user not found",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user by setting userID in context
				ctx.Set("userID", "nonexistent-user-id")
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("GetUserByID", mock.Anything, "nonexistent-user-id").
					Return(nil, dbutils.ErrRecordNotFoundType)
				return mockUserSvc
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "service layer error",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user by setting userID in context
				ctx.Set("userID", "de305d54-75b4-431b-adb2-eb6b9e546099")
			},

			setupMockSvc: func(ctx *gin.Context) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("GetUserByID", mock.Anything, "de305d54-75b4-431b-adb2-eb6b9e546099").
					Return(nil, assert.AnError)
				return mockUserSvc
			},

			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"message":"Internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx)
			mockUserSvc := tc.setupMockSvc(ctx)

			userHandler := NewUser(mockUserSvc)
			userHandler.GetProfile(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestUser_UpdateProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputRequest *updateProfileRequest

		setupRequest func(ctx *gin.Context, inputRequest *updateProfileRequest)

		setupMockSvc func(ctx *gin.Context, inputRequest *updateProfileRequest) *svcMocks.User

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful update profile",

			inputRequest: &updateProfileRequest{
				DisplayName: "Updated User",
				Email:       "updateduser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *updateProfileRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user by setting userID in context
				ctx.Set("userID", "de305d54-75b4-431b-adb2-eb6b9e546099")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *updateProfileRequest) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("UpdateUserByID", ctx, "de305d54-75b4-431b-adb2-eb6b9e546099", inputRequest.DisplayName, inputRequest.Email).
					Return(nil)
				return mockUserSvc
			},

			expectedCode:     http.StatusOK,
			expectedResponse: `{"message":"Edit current user successfully!"}`,
		},
		{
			name: "unauthenticated request",

			inputRequest: &updateProfileRequest{
				DisplayName: "Updated User",
				Email:       "updateduser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *updateProfileRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// No userID set in context to simulate unauthenticated request
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *updateProfileRequest) *svcMocks.User {
				return svcMocks.NewUser(t) // No expectations since service should not be called
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "invalid userID in context",

			inputRequest: &updateProfileRequest{
				DisplayName: "Updated User",
				Email:       "updateduser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *updateProfileRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Set invalid userID in context
				ctx.Set("userID", 12345) // should be a string
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *updateProfileRequest) *svcMocks.User {
				return svcMocks.NewUser(t) // No expectations since service should not be called
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "user not found",

			inputRequest: &updateProfileRequest{
				DisplayName: "Updated User",
				Email:       "updateduser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *updateProfileRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user by setting userID in context
				ctx.Set("userID", "nonexistent-user-id")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *updateProfileRequest) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("UpdateUserByID", ctx, "nonexistent-user-id", inputRequest.DisplayName, inputRequest.Email).
					Return(dbutils.ErrRecordNotFoundType)
				return mockUserSvc
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "service layer error",

			inputRequest: &updateProfileRequest{
				DisplayName: "Updated User",
				Email:       "updateduser@example.com",
			},

			setupRequest: func(ctx *gin.Context, inputRequest *updateProfileRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
				// Simulate authenticated user by setting userID in context
				ctx.Set("userID", "de305d54-75b4-431b-adb2-eb6b9e546099")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *updateProfileRequest) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("UpdateUserByID", ctx, "de305d54-75b4-431b-adb2-eb6b9e546099", inputRequest.DisplayName, inputRequest.Email).
					Return(assert.AnError)
				return mockUserSvc
			},

			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"message":"Internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx, tc.inputRequest)
			mockUserSvc := tc.setupMockSvc(ctx, tc.inputRequest)

			userHandler := NewUser(mockUserSvc)
			userHandler.UpdateProfile(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}
