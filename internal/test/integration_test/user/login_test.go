package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
	"github.com/vukieuhaihoa/bookmark-management/internal/test/fixture"
	"github.com/vukieuhaihoa/bookmark-management/pkg/jwtutils/mocks"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-management/pkg/utils"
)

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
			apiEngine := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark_service",
					InstanceID:  "test_instance_id_1",
				},
				RedisClient:     redisPkg.InitMockRedis(t),
				SqlDB:           db,
				RandomCodeGen:   nil,
				PasswordHashing: utils.NewPasswordHashing(),
				JWTGenerator:    jwtGen,
				JWTValidator:    nil,
			})

			// Setup test HTTP request
			respRec := tc.setupTestHTTP(apiEngine)

			// Verify response status code
			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			assert.Contains(t, respRec.Body.String(), tc.expectedMessageResponse)
		})
	}
}
