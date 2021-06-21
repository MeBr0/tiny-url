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
	"net/http/httptest"
	"testing"
)

func TestHandler_register(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAuth, user domain.UserRegister)

	tests := []struct {
		name          string
		requestBody   string
		requestUser   domain.UserRegister
		mockBehaviour mockBehaviour
		statusCode    int
		responseBody  string
	}{
		{
			name:        "ok",
			requestBody: `{"name": "mebr", "email": "qweqweqwe@gmail.com", "password": "qweqweqwe"}`,
			requestUser: domain.UserRegister{
				Name:     "mebr",
				Email:    "qweqweqwe@gmail.com",
				Password: "qweqweqwe",
			},
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserRegister) {
				s.EXPECT().Register(context.Background(), user).Return(nil)
			},
			statusCode:   201,
			responseBody: ``,
		},
		{
			name:          "invalid request body",
			requestBody:   `{}`,
			requestUser:   domain.UserRegister{},
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserRegister) {},
			statusCode:    400,
			responseBody:  `{"message":"invalid request body"}`,
		},
		{
			name:        "user already exists",
			requestBody: `{"name": "mebr", "email": "qweqweqwe@gmail.com", "password": "qweqweqwe"}`,
			requestUser: domain.UserRegister{
				Name:     "mebr",
				Email:    "qweqweqwe@gmail.com",
				Password: "qweqweqwe",
			},
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserRegister) {
				s.EXPECT().Register(context.Background(), user).Return(repo.ErrUserAlreadyExists)
			},
			statusCode:   400,
			responseBody: `{"message":"user already exists"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuth(c)
			tt.mockBehaviour(auth, tt.requestUser)

			services := &service.Services{Auth: auth}
			handler := &Handler{
				services:     services,
				tokenManager: nil,
			}

			// Init Endpoint
			r := gin.New()
			r.POST("/register", handler.register)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(tt.requestBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_login(t *testing.T) {
	type mockBehaviour func(s *mockService.MockAuth, user domain.UserLogin)

	tests := []struct {
		name          string
		requestBody   string
		requestUser   domain.UserLogin
		mockBehaviour mockBehaviour
		statusCode    int
		responseBody  string
	}{
		{
			name:        "ok",
			requestBody: `{"email": "qweqweqwe@gmail.com", "password": "qweqweqwe"}`,
			requestUser: domain.UserLogin{
				Email:    "qweqweqwe@gmail.com",
				Password: "qweqweqwe",
			},
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserLogin) {
				s.EXPECT().Login(context.Background(), user).Return(domain.Tokens{
					AccessToken: "token",
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"accessToken":"token"}`,
		},
		{
			name:          "invalid request body",
			requestBody:   `{"email": "qweqweqwe", "password": "qweqweqwe"}`,
			requestUser:   domain.UserLogin{},
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserLogin) {},
			statusCode:    400,
			responseBody:  `{"message":"invalid request body"}`,
		},
		{
			name:        "user does not exists",
			requestBody: `{"email": "qweqweqwe@gmail.com", "password": "qweqweqwe"}`,
			requestUser: domain.UserLogin{
				Email:    "qweqweqwe@gmail.com",
				Password: "qweqweqwe",
			},
			mockBehaviour: func(s *mockService.MockAuth, user domain.UserLogin) {
				s.EXPECT().Login(context.Background(), user).Return(domain.Tokens{}, repo.ErrUserNotFound)
			},
			statusCode:   400,
			responseBody: `{"message":"user doesn't exists"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuth(c)
			tt.mockBehaviour(auth, tt.requestUser)

			services := &service.Services{Auth: auth}
			handler := &Handler{
				services:     services,
				tokenManager: nil,
			}

			// Init Endpoint
			r := gin.New()
			r.POST("/login", handler.login)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(tt.requestBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}
