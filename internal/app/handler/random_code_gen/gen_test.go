package random_code_gen

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-management/internal/app/service/random_code_gen/mocks"
)

func TestHandler_GeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func() *mocks.Service

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "successful password generation",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/generate-password", nil)
			},
			setupMockSvc: func() *mocks.Service {
				svcMock := mocks.NewService(t)
				svcMock.On("GeneratePassword").Return("securepassword123", nil)
				return svcMock
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "securepassword123",
		},
		{
			name: "password generation failure",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/generate-password", nil)
			},
			setupMockSvc: func() *mocks.Service {
				svcMock := mocks.NewService(t)
				svcMock.On("GeneratePassword").Return("", errors.New("something went wrong"))
				return svcMock
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: ErrPasswordGenerationFailed.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)

			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc()
			testHandler := NewPasswordHandler(mockSvc)
			testHandler.GeneratePassword(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
