package v1

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/internal/repo"
	"github.com/mebr0/tiny-url/internal/service"
	mockService "github.com/mebr0/tiny-url/internal/service/mocks"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_redirectWithAlias(t *testing.T) {
	type mockBehaviour func(s *mockService.MockURLs, alias string)

	userId := primitive.NewObjectID()

	tests := []struct {
		name          string
		alias         string
		mockBehaviour mockBehaviour
		statusCode    int
		responseBody  string
	}{
		{
			name:  "ok",
			alias: "alias",
			mockBehaviour: func(s *mockService.MockURLs, alias string) {
				s.EXPECT().Get(context.Background(), alias).Return(domain.URL{
					Alias:     "alias",
					Original:  "https://google.com",
					CreatedAt: time.Now(),
					ExpiredAt: time.Now().Add(5 * time.Minute),
					Owner:     userId,
				}, nil)
			},
			statusCode:   301,
			responseBody: ``,
		},
		{
			name:  "url expired",
			alias: "alias",
			mockBehaviour: func(s *mockService.MockURLs, alias string) {
				s.EXPECT().Get(context.Background(), alias).Return(domain.URL{
					Alias:     "alias",
					Original:  "https://google.com",
					CreatedAt: time.Now(),
					ExpiredAt: time.Now().Add(-5 * time.Minute),
					Owner:     userId,
				}, nil)
			},
			statusCode:   400,
			responseBody: `{"message":"url expired"}`,
		},
		{
			name:  "url does not exists",
			alias: "alias",
			mockBehaviour: func(s *mockService.MockURLs, alias string) {
				s.EXPECT().Get(context.Background(), alias).Return(domain.URL{}, repo.ErrURLNotFound)
			},
			statusCode:   400,
			responseBody: `{"message":"url doesn't exists"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			urlsService := mockService.NewMockURLs(c)
			tt.mockBehaviour(urlsService, tt.alias)

			services := &service.Services{URLs: urlsService}
			handler := &Handler{
				services:     services,
				tokenManager: nil,
			}

			// Init Endpoint
			r := gin.New()
			r.GET("/to/:alias", handler.redirectWithAlias)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/to/alias", bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)

			if tt.responseBody != "" {
				assert.Equal(t, tt.responseBody, w.Body.String())
			}
		})
	}
}
