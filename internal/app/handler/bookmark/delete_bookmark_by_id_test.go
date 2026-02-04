package bookmark

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/service/bookmark/mocks"
	"github.com/vukieuhaihoa/bookmark-management/pkg/dbutils"
)

func TestHandler_DeleteBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputBookmarkID string
		inputUserID     string

		setupMockRequest func(c *gin.Context, id, userID string)

		setupMockBookmarkService func(ctx context.Context, id, userID string) *mocks.Service

		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Delete bookmark by ID successfully",

			setupMockRequest: func(c *gin.Context, id, userID string) {
				c.Request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/bookmarks/%s", id), nil)
				c.Params = gin.Params{{Key: "id", Value: id}}
				// Simulate authenticated user
				c.Set("claims", jwt.MapClaims{"sub": userID})
			},

			setupMockBookmarkService: func(ctx context.Context, id, userID string) *mocks.Service {
				repoMock := mocks.NewService(t)
				repoMock.On("DeleteBookmarkByID", ctx, id, userID).Return(nil)
				return repoMock
			},
			inputBookmarkID:    "a1b2c3d4-e5f6-7890-abcd-ef0000000006",
			inputUserID:        "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"message":"Success"}`,
		},
		{
			name: "Fail to delete bookmark by ID - bookmark not found",

			setupMockRequest: func(c *gin.Context, id, userID string) {
				c.Request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/bookmarks/%s", id), nil)
				c.Params = gin.Params{{Key: "id", Value: id}}
				// Simulate authenticated user
				c.Set("claims", jwt.MapClaims{"sub": userID})
			},

			setupMockBookmarkService: func(ctx context.Context, id, userID string) *mocks.Service {
				repoMock := mocks.NewService(t)
				repoMock.On("DeleteBookmarkByID", ctx, id, userID).Return(dbutils.ErrRecordNotFoundType)
				return repoMock
			},
			inputBookmarkID:    "non-existent-id",
			inputUserID:        "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"Invalid input"}`,
		},
		{
			name: "Fail to delete bookmark by ID - invalid user id",

			setupMockRequest: func(c *gin.Context, id, userID string) {
				c.Request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/bookmarks/%s", id), nil)
				c.Params = gin.Params{{Key: "id", Value: id}}
				// Simulate authenticated user
				c.Set("claims", jwt.MapClaims{"sub": userID})
			},

			setupMockBookmarkService: func(ctx context.Context, id, userID string) *mocks.Service {
				repoMock := mocks.NewService(t)
				repoMock.On("DeleteBookmarkByID", ctx, id, userID).Return(dbutils.ErrRecordNotFoundType)
				return repoMock
			},
			inputBookmarkID:    "a1b2c3d4-e5f6-7890-abcd-ef0000000006",
			inputUserID:        "invalid-user-id",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"Invalid input"}`,
		},
		{
			name: "Fail to delete bookmark by ID - internal server error",

			setupMockRequest: func(c *gin.Context, id, userID string) {
				c.Request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/bookmarks/%s", id), nil)
				c.Params = gin.Params{{Key: "id", Value: id}}
				// Simulate authenticated user
				c.Set("claims", jwt.MapClaims{"sub": userID})
			},

			setupMockBookmarkService: func(ctx context.Context, id, userID string) *mocks.Service {
				repoMock := mocks.NewService(t)
				repoMock.On("DeleteBookmarkByID", ctx, id, userID).Return(fmt.Errorf("database error"))
				return repoMock
			},
			inputBookmarkID:    "a1b2c3d4-e5f6-7890-abcd-ef0000000006",
			inputUserID:        "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"message":"Internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			// Setup mock request
			tc.setupMockRequest(ctx, tc.inputBookmarkID, tc.inputUserID)

			// Setup mock bookmark service
			bookmarkServiceMock := tc.setupMockBookmarkService(ctx, tc.inputBookmarkID, tc.inputUserID)

			// Create handler with mock service
			handler := NewBookmarkHandler(bookmarkServiceMock)

			// Call DeleteBookmarkByID handler
			handler.DeleteBookmarkByID(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
