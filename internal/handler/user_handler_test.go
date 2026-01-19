package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	svcMocks "github.com/vukieuhaihoa/bookmark-management/internal/service/mocks"
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

		expectedCode   int
		verifyResponse func(t *testing.T, rec *httptest.ResponseRecorder, inputRequest *createUserRequest)
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
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/user/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("CreateUser", mock.Anything, inputRequest.Username, inputRequest.Password, inputRequest.DisplayName, inputRequest.Email).
					Return(&model.User{
						Username:    inputRequest.Username,
						DisplayName: inputRequest.DisplayName,
						Email:       inputRequest.Email,
						ID:          "de305d54-75b4-431b-adb2-eb6b9e546099",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}, nil)
				return mockUserSvc
			},

			expectedCode: http.StatusCreated,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder, inputRequest *createUserRequest) {
				// Parse the response body
				var response createUserResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err, "Response should be valid JSON")

				// Verify response structure
				assert.NotNil(t, response.Data, "Data should not be nil")
				assert.Equal(t, "Register an user successfully!", response.Message)

				// Verify user data
				assert.Equal(t, "de305d54-75b4-431b-adb2-eb6b9e546099", response.Data.ID)
				assert.Equal(t, inputRequest.Username, response.Data.Username)
				assert.Equal(t, inputRequest.DisplayName, response.Data.DisplayName)
				assert.Equal(t, inputRequest.Email, response.Data.Email)

				// Verify timestamps are set
				assert.False(t, response.Data.CreatedAt.IsZero(), "CreatedAt should be set")
				assert.False(t, response.Data.UpdatedAt.IsZero(), "UpdatedAt should be set")

				// Password should not be in the response (due to json:"-" tag)
				assert.Empty(t, response.Data.Password, "Password should not be in response")
			},
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
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/user/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.User {
				return svcMocks.NewUser(t) // No expectations since service should not be called
			},

			expectedCode: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder, inputRequest *createUserRequest) {
				// Parse the response body
				var respBody map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &respBody)
				assert.NoError(t, err, "Response should be valid JSON")

				// Verify error message structure
				assert.Contains(t, respBody, "details")
				details, ok := respBody["details"].([]interface{})
				assert.True(t, ok, "Error details should be a list")

				expectedErrors := []string{
					"Username is invalid (required)",
					"Password is invalid (min)",
					"Email is invalid (email)",
				}

				var detailStrs []string
				for _, d := range details {
					if ds, ok := d.(string); ok {
						detailStrs = append(detailStrs, ds)
					}
				}

				for _, expectedErr := range expectedErrors {
					assert.Contains(t, detailStrs, expectedErr)
				}
			},
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
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/user/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.User {
				return svcMocks.NewUser(t) // No expectations since service should not be called
			},

			expectedCode: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder, inputRequest *createUserRequest) {
				// Parse the response body
				var respBody map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &respBody)
				assert.NoError(t, err, "Response should be valid JSON")

				// Verify error message structure
				assert.Contains(t, respBody, "details")
				details, ok := respBody["details"].([]interface{})
				assert.True(t, ok, "Error details should be a list")

				expectedErrors := []string{
					"Password is invalid (password_strength)",
				}

				var detailStrs []string
				for _, d := range details {
					if ds, ok := d.(string); ok {
						detailStrs = append(detailStrs, ds)
					}
				}

				for _, expectedErr := range expectedErrors {
					assert.Contains(t, detailStrs, expectedErr)
				}
			},
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
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/user/register", strings.NewReader(string(reqBody)))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			setupMockSvc: func(ctx *gin.Context, inputRequest *createUserRequest) *svcMocks.User {
				mockUserSvc := svcMocks.NewUser(t)
				mockUserSvc.On("CreateUser", mock.Anything, inputRequest.Username, inputRequest.Password, inputRequest.DisplayName, inputRequest.Email).
					Return(nil, assert.AnError)
				return mockUserSvc
			},

			expectedCode: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder, inputRequest *createUserRequest) {
				// Parse the response body
				var respBody map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &respBody)
				assert.NoError(t, err, "Response should be valid JSON")

				// Verify error message
				assert.Equal(t, "Internal server error", respBody["message"])
			},
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
			if tc.verifyResponse != nil {
				tc.verifyResponse(t, rec, tc.inputRequest)
			}
		})
	}
}
