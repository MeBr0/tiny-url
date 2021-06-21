package v1

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestHandler_listURLs(t *testing.T) {
	type mockBehaviour func(s *mockService.MockURLs, ownerId primitive.ObjectID)

	userId := primitive.NewObjectID()

	urls := []domain.URL{
		{
			Alias:     "alias",
			Original:  "https://google.com",
			CreatedAt: time.Now(),
			ExpiredAt: time.Now(),
			Owner:     userId,
		},
	}

	setResponseBody := func(urls []domain.URL) string {
		body, _ := json.Marshal(urls)

		return string(body)
	}

	tests := []struct {
		name          string
		userId        primitive.ObjectID
		mockBehaviour mockBehaviour
		statusCode    int
		responseBody  string
	}{
		{
			name:   "ok",
			userId: userId,
			mockBehaviour: func(s *mockService.MockURLs, ownerId primitive.ObjectID) {
				s.EXPECT().ListByOwner(context.Background(), ownerId).Return(urls, nil)
			},
			statusCode:   200,
			responseBody: setResponseBody(urls),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			urlsService := mockService.NewMockURLs(c)
			tt.mockBehaviour(urlsService, tt.userId)

			services := &service.Services{URLs: urlsService}
			handler := &Handler{
				services:     services,
				tokenManager: nil,
			}

			// Init Endpoint
			r := gin.New()
			r.GET("/urls", func(c *gin.Context) {
				c.Set(userCtx, userId.Hex())
			}, handler.listURLs)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/urls", bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_createURL(t *testing.T) {
	type mockBehaviour func(s *mockService.MockURLs, url domain.URLCreate)

	userId := primitive.NewObjectID()

	tests := []struct {
		name          string
		userId        primitive.ObjectID
		requestBody   string
		requestURL    domain.URLCreate
		mockBehaviour mockBehaviour
		statusCode    int
		responseBody  string
	}{
		{
			name:        "ok",
			userId:      userId,
			requestBody: `{"original": "https://google.com", "duration": 60}`,
			requestURL: domain.URLCreate{
				Original: "https://google.com",
				Duration: 60,
			},
			mockBehaviour: func(s *mockService.MockURLs, url domain.URLCreate) {
				toCreate := domain.URLCreate{
					Original: url.Original,
					Duration: url.Duration,
					Owner:    userId,
				}

				created := domain.NewURL(toCreate, "alias")

				s.EXPECT().Create(context.Background(), toCreate).Return(created, nil)
			},
			statusCode:   201,
			responseBody: ``,
		},
		{
			name:          "invalid request body",
			requestBody:   `{"duration": 60}`,
			mockBehaviour: func(s *mockService.MockURLs, url domain.URLCreate) {},
			statusCode:    400,
			responseBody:  `{"message":"invalid request body"}`,
		},
		{
			name:        "url already exists",
			userId:      userId,
			requestBody: `{"original": "https://google.com", "duration": 60}`,
			requestURL: domain.URLCreate{
				Original: "https://google.com",
				Duration: 60,
			},
			mockBehaviour: func(s *mockService.MockURLs, url domain.URLCreate) {
				toCreate := domain.URLCreate{
					Original: url.Original,
					Duration: url.Duration,
					Owner:    userId,
				}

				s.EXPECT().Create(context.Background(), toCreate).Return(domain.URL{}, repo.ErrURLAlreadyExists)
			},
			statusCode:   400,
			responseBody: `{"message":"url already exists"}`,
		},
		{
			name:        "url limit",
			userId:      userId,
			requestBody: `{"original": "https://google.com", "duration": 60}`,
			requestURL: domain.URLCreate{
				Original: "https://google.com",
				Duration: 60,
			},
			mockBehaviour: func(s *mockService.MockURLs, url domain.URLCreate) {
				toCreate := domain.URLCreate{
					Original: url.Original,
					Duration: url.Duration,
					Owner:    userId,
				}

				s.EXPECT().Create(context.Background(), toCreate).Return(domain.URL{}, service.ErrURLLimit)
			},
			statusCode:   400,
			responseBody: `{"message":"cannot create more urls"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			urlsService := mockService.NewMockURLs(c)
			tt.mockBehaviour(urlsService, tt.requestURL)

			services := &service.Services{URLs: urlsService}
			handler := &Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.POST("/urls", func(c *gin.Context) {
				c.Set(userCtx, userId.Hex())
			}, handler.createURL)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/urls", bytes.NewBufferString(tt.requestBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)

			// Since created URL can have arbitrary dates, not assert empty body
			if tt.responseBody != "" {
				assert.Equal(t, tt.responseBody, w.Body.String())
			}
		})
	}
}
