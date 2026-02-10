package link

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/api"
	redisPkg "github.com/vukieuhaihoa/bookmark-management/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-management/pkg/sqldb"
	"github.com/vukieuhaihoa/bookmark-management/pkg/utils"
)

func TestShortenURLEndpoint_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode int
		expectedCodeLength int
	}{
		{
			name: "successful shorten URL",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/links/shorten", strings.NewReader(`{"url":"http://example.com","exp":3600}`)) // Body would be added
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusOK,
			expectedCodeLength: 10,
		},
		{
			name: "invalid request payload",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/links/shorten", strings.NewReader(`{"url":"", "exp":-1}`)) // Body would be added
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedCodeLength: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			apiEngine := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark-service",
					InstanceID:  "test_instance_id_1",
				},
				RedisClient:     redisPkg.InitMockRedis(t),
				SqlDB:           sqldb.InitMockDB(t),
				RandomCodeGen:   utils.NewCodeGenerator(),
				PasswordHashing: nil,
				JWTGenerator:    nil,
				JWTValidator:    nil,
			})

			respRec := tc.setupTestHTTP(apiEngine)

			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			// Check code length
			var respBody struct {
				Code string `json:"code"`
			}
			err := json.Unmarshal(respRec.Body.Bytes(), &respBody)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCodeLength, len(respBody.Code))
		})
	}
}
