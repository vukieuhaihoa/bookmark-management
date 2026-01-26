package endpoint

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
	"github.com/vukieuhaihoa/bookmark-management/internal/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/jwtutils/mocks"
	"github.com/vukieuhaihoa/bookmark-management/pkg/sqldb"
)

func TestUserEndpoint_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode      int
		expectedMessageResponse string
	}{
		{
			name: "successful user registration",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(`{"username":"testuser","password":"my_SECURE_password123@","display_name":"Test User","email":"testuser@example.com"}`))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusCreated,
			expectedMessageResponse: `message":"Register an user successfully!"`,
		},
		{
			name: "invalid user registration payload",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(`{"username":"","password":"weak","display_name":"Test User","email":"invalid-email"}`))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusBadRequest,
			expectedMessageResponse: `"message":"Invalid input fields"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// init mock db and migrate
			db := sqldb.InitMockDB(t)
			db.AutoMigrate(&model.User{})

			// Initialize API engine
			apiEngine := api.New(&api.Config{
				ServiceName: "bookmark_service",
				InstanceID:  "test_instance_id_1",
			}, nil, db, nil, nil)

			// Setup test HTTP request
			respRec := tc.setupTestHTTP(apiEngine)

			// Verify response status code
			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			assert.Contains(t, respRec.Body.String(), tc.expectedMessageResponse)
		})
	}
}

func TestUserEndpoint_Login(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		setupMockJWTGenerator func(t *testing.T) *mocks.JWTGenerator

		expectedStatusCode      int
		expectedMessageResponse string
	}{
		{
			name: "successful user login",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(`{"username":"testuser001","password":"my_SECURE_password123@"}`))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTGenerator: func(t *testing.T) *mocks.JWTGenerator {
				jwtGen := mocks.NewJWTGenerator(t)
				jwtGen.On("GenerateToken", mock.Anything).Return("mocked_jwt_token", nil)
				return jwtGen
			},

			expectedStatusCode:      http.StatusOK,
			expectedMessageResponse: `"message":"Logged in successfully!"`,
		},
		{
			name: "invalid user login payload",

			setupMockJWTGenerator: func(t *testing.T) *mocks.JWTGenerator {
				return mocks.NewJWTGenerator(t)
			},

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(`{"username":"","password":""}`))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusBadRequest,
			expectedMessageResponse: `"message":"Invalid input fields"`,
		},
		{
			name: "user login failed - invalid credentials",

			setupMockJWTGenerator: func(t *testing.T) *mocks.JWTGenerator {
				return mocks.NewJWTGenerator(t)
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(`{"username":"testuser001","password":"wrong_password"}`))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusBadRequest,
			expectedMessageResponse: `"message":"invalid username or password"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// init mock db and migrate
			db := fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			jwtGen := tc.setupMockJWTGenerator(t)

			// Initialize API engine
			apiEngine := api.New(&api.Config{
				ServiceName: "bookmark_service",
				InstanceID:  "test_instance_id_1",
			}, nil, db, jwtGen, nil)

			// Setup test HTTP request
			respRec := tc.setupTestHTTP(apiEngine)

			// Verify response status code
			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			assert.Contains(t, respRec.Body.String(), tc.expectedMessageResponse)
		})
	}
}

func TestUserEndpoint_GetProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		setupMockJWTValidator func(t *testing.T) *mocks.JWTValidator

		expectedStatusCode      int
		expectedMessageResponse string
	}{
		{
			name: "successful get user profile",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("GET", "/v1/self/info", nil)
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedStatusCode:      http.StatusOK,
			expectedMessageResponse: `"username":"testuser001"`,
		},
		{
			name: "get user profile failed - invalid token",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("GET", "/v1/self/info", nil)
				req.Header.Set("Authorization", "Bearer invalid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "invalid_jwt_token").Return(nil, assert.AnError)
				return jwtValidator
			},

			expectedStatusCode:      http.StatusUnauthorized,
			expectedMessageResponse: `"message":"Invalid token"`,
		},
		{
			name: "token does not contain user ID",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("GET", "/v1/self/info", nil)
				req.Header.Set("Authorization", "Bearer token_without_user_id")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "token_without_user_id").Return(jwt.MapClaims{}, nil)
				return jwtValidator
			},

			expectedStatusCode:      http.StatusUnauthorized,
			expectedMessageResponse: `"message":"Invalid token"`,
		},
		{
			name: "user not found",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("GET", "/v1/self/info", nil)
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "non_existent_user_id"}, nil)
				return jwtValidator
			},

			expectedStatusCode:      http.StatusUnauthorized,
			expectedMessageResponse: `"message":"Unauthorized"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// init mock db and migrate
			db := fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			jwtValidator := tc.setupMockJWTValidator(t)

			// Initialize API engine
			apiEngine := api.New(&api.Config{
				ServiceName: "bookmark_service",
				InstanceID:  "test_instance_id_1",
			}, nil, db, nil, jwtValidator)

			// Setup test HTTP request
			respRec := tc.setupTestHTTP(apiEngine)

			// Verify response status code
			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			assert.Contains(t, respRec.Body.String(), tc.expectedMessageResponse)
		})
	}
}

func TestUserEndpoint_UpdateProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		setupMockJWTValidator func(t *testing.T) *mocks.JWTValidator

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful update user profile",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("PUT", "/v1/self/info", strings.NewReader(`{"display_name":"Test User 1 Updated","email":"testuser001updated@example.com"}`))
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedCode:     http.StatusOK,
			expectedResponse: `"message":"Edit current user successfully!"`,
		},
		{
			name: "get user profile failed - invalid token",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("PUT", "/v1/self/info", strings.NewReader(`{"display_name":"Test User 1 Updated","email":"testuser001updated@example.com"`))
				req.Header.Set("Authorization", "Bearer invalid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "invalid_jwt_token").Return(nil, assert.AnError)
				return jwtValidator
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `"message":"Invalid token"`,
		},
		{
			name: "token does not contain user ID",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("PUT", "/v1/self/info", strings.NewReader(`{"display_name":"Test User 1 Updated","email":"testuser001updated@example.com"`))
				req.Header.Set("Authorization", "Bearer token_without_user_id")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "token_without_user_id").Return(jwt.MapClaims{}, nil)
				return jwtValidator
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `"message":"Invalid token"`,
		},
		{
			name: "user not found",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("PUT", "/v1/self/info", strings.NewReader(`{"display_name":"Test User 1 Updated","email":"testuser001updated@example.com"}`))
				req.Header.Set("Authorization", "Bearer user_not_found_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "user_not_found_token").Return(jwt.MapClaims{"sub": "non_existent_user_id"}, nil)
				return jwtValidator
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `"message":"Unauthorized"`,
		},
		{
			name: "update user profile failed - invalid payload",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("PUT", "/v1/self/info", strings.NewReader(`{"display_name":"","email":"invalid-email"}`))
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `"message":"Invalid input fields"`,
		},
		{
			name: "update user profile failed - email already exists",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("PUT", "/v1/self/info", strings.NewReader(`{"display_name":"Test User 1 Updated","email":"alice@example.com"}`))
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `"message":"email already exists"`,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			jwtValidator := tc.setupMockJWTValidator(t)

			// Initialize API engine
			apiEngine := api.New(&api.Config{
				ServiceName: "bookmark_service",
				InstanceID:  "test_instance_id_1",
			}, nil, db, nil, jwtValidator)

			// Setup test HTTP request
			respRec := tc.setupTestHTTP(apiEngine)

			// Verify response status code
			assert.Equal(t, tc.expectedCode, respRec.Code)
			assert.Contains(t, respRec.Body.String(), tc.expectedResponse)
		})
	}
}
