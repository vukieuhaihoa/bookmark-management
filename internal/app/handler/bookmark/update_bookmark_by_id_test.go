package bookmark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/service/bookmark/mocks"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
)

func TestHandler_UpdateBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputBookmarkID string
		inputUserID     string
		inputRequest    *updateBookmarkRequest

		setupMockRequest func(c *gin.Context, id, userID string, inputRequest *updateBookmarkRequest)

		setupMockBookmarkService func(ctx context.Context, id, userID string, req *updateBookmarkRequest) *mocks.Service

		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Update bookmark by ID successfully",

			setupMockRequest: func(c *gin.Context, id, userID string, inputRequest *updateBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				c.Request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/v1/bookmarks/%s", id), strings.NewReader(string(reqBody)))
				c.Request.Header.Set("Content-Type", "application/json")
				c.Params = gin.Params{{Key: "id", Value: id}}
				// Simulate authenticated user
				c.Set("claims", jwt.MapClaims{"sub": userID})
			},

			setupMockBookmarkService: func(ctx context.Context, id, userID string, req *updateBookmarkRequest) *mocks.Service {
				repoMock := mocks.NewService(t)
				repoMock.On("UpdateBookmarkByID", ctx, id, userID, &model.Bookmark{
					URL:         req.URL,
					Description: req.Description,
				}).Return(nil)
				return repoMock
			},
			inputBookmarkID: "a1b2c3d4-e5f6-7890-abcd-ef0000000005",
			inputUserID:     "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			inputRequest: &updateBookmarkRequest{
				URL:         "https://updated-example.com",
				Description: "This is an updated description.",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"message":"Success"}`,
		},
		{
			name: "Fail to update bookmark by ID - bookmark not found",

			setupMockRequest: func(c *gin.Context, id, userID string, inputRequest *updateBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				c.Request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/v1/bookmarks/%s", id), strings.NewReader(string(reqBody)))
				c.Request.Header.Set("Content-Type", "application/json")
				c.Params = gin.Params{{Key: "id", Value: id}}
				// Simulate authenticated user
				c.Set("claims", jwt.MapClaims{"sub": userID})
			},

			setupMockBookmarkService: func(ctx context.Context, id, userID string, req *updateBookmarkRequest) *mocks.Service {
				repoMock := mocks.NewService(t)
				repoMock.On("UpdateBookmarkByID", ctx, id, userID, &model.Bookmark{
					URL:         req.URL,
					Description: req.Description,
				}).Return(dbutils.ErrRecordNotFoundType)
				return repoMock
			},
			inputBookmarkID: "non-existent-id",
			inputUserID:     "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			inputRequest: &updateBookmarkRequest{
				URL:         "https://nonexistent.com",
				Description: "This bookmark does not exist.",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "{\"message\":\"Invalid input\"}",
		},
		{
			name: "Fail to update bookmark by ID - invalid user ID",

			setupMockRequest: func(c *gin.Context, id, userID string, inputRequest *updateBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				c.Request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/v1/bookmarks/%s", id), strings.NewReader(string(reqBody)))
				c.Request.Header.Set("Content-Type", "application/json")
				c.Params = gin.Params{{Key: "id", Value: id}}
				// Simulate authenticated user
				c.Set("claims", jwt.MapClaims{"sub": userID})
			},

			setupMockBookmarkService: func(ctx context.Context, id, userID string, req *updateBookmarkRequest) *mocks.Service {
				repoMock := mocks.NewService(t)
				repoMock.On("UpdateBookmarkByID", ctx, id, userID, &model.Bookmark{
					URL:         req.URL,
					Description: req.Description,
				}).Return(dbutils.ErrRecordNotFoundType)
				return repoMock
			},
			inputBookmarkID: "a1b2c3d4-e5f6-7890-abcd-ef0000000005",
			inputUserID:     "invalid-user-id",
			inputRequest: &updateBookmarkRequest{
				URL:         "https://updated-example.com",
				Description: "This is an updated description.",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "{\"message\":\"Invalid input\"}",
		},
		{
			name: "Fail to update bookmark by ID - internal server error",

			setupMockRequest: func(c *gin.Context, id, userID string, inputRequest *updateBookmarkRequest) {
				reqBody, _ := json.Marshal(inputRequest)
				c.Request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/v1/bookmarks/%s", id), strings.NewReader(string(reqBody)))
				c.Request.Header.Set("Content-Type", "application/json")
				c.Params = gin.Params{{Key: "id", Value: id}}
				// Simulate authenticated user
				c.Set("claims", jwt.MapClaims{"sub": userID})
			},
			setupMockBookmarkService: func(ctx context.Context, id, userID string, req *updateBookmarkRequest) *mocks.Service {
				repoMock := mocks.NewService(t)
				repoMock.On("UpdateBookmarkByID", ctx, id, userID, &model.Bookmark{
					URL:         req.URL,
					Description: req.Description,
				}).Return(assert.AnError)
				return repoMock
			},

			inputBookmarkID: "a1b2c3d4-e5f6-7890-abcd-ef0000000005",
			inputUserID:     "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			inputRequest: &updateBookmarkRequest{
				URL:         "https://error-example.com",
				Description: "This will cause an internal error.",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "{\"message\":\"Internal server error\"}",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup Gin context with request
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupMockRequest(ctx, tc.inputBookmarkID, tc.inputUserID, tc.inputRequest)

			// Setup mock service
			mockSvc := tc.setupMockBookmarkService(ctx, tc.inputBookmarkID, tc.inputUserID, tc.inputRequest)

			// Create handler with mock service
			handler := NewBookmarkHandler(mockSvc)

			// Call the handler
			handler.UpdateBookmarkByID(ctx)

			// Assert response
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}
