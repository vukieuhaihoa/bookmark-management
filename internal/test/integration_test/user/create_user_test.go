package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/pkg/sqldb"
	"github.com/vukieuhaihoa/bookmark-management/pkg/utils"
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
			apiEngine := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark_service",
					InstanceID:  "test_instance_id_1",
				},
				RedisClient:     nil,
				SqlDB:           db,
				RandomCodeGen:   nil,
				PasswordHashing: utils.NewPasswordHashing(),
				JWTGenerator:    nil,
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
