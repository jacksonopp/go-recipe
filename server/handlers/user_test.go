package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jacksonopp/go-recipe/domain"
	"github.com/valyala/fasthttp"
	"testing"
)

type mockUserService struct{}

func (s *mockUserService) CreateUser(user domain.User) error {
	return nil
}

func (s *mockUserService) GetUserByName(name string) (*domain.User, error) {
	return nil, nil
}

func (s *mockUserService) LoginUser(name, password string) (*domain.User, error) {
	return nil, nil
}

func TestUserHandler_register(t *testing.T) {
	app := fiber.New()
	h := UserHandler{s: &mockUserService{}, r: app}
	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "valid request",
			body:           `{"username": "test", "password": "test", "passwordConfirm": "test"}`,
			expectedStatus: 201,
		},
		{
			name:           "invalid request",
			body:           `{`,
			expectedStatus: fiber.StatusUnprocessableEntity,
		},
		{
			name:           "data doesnt match",
			body:           `{"badData": "badData"}`,
			expectedStatus: fiber.StatusUnprocessableEntity,
		},
		{
			name:           "missing username",
			body:           `{"username": "", password": "test", "passwordConfirm": "test"}`,
			expectedStatus: fiber.StatusUnprocessableEntity,
		},
		{
			name:           "missing password",
			body:           `{"username": "test", password": "", "passwordConfirm": ""}`,
			expectedStatus: fiber.StatusUnprocessableEntity,
		},
		{
			name:           "passwords dont match",
			body:           `{"username": "test", password": "test", "passwordConfirm": "test2"}`,
			expectedStatus: fiber.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Request().SetBodyString(tt.body)

			//set the request to application/json
			c.Request().Header.SetMethod("POST")
			c.Request().Header.SetContentType("application/json")

			_ = h.register(c)

			if c.Response().StatusCode() != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, c.Response().StatusCode())
			}
		})
	}
}

func TestUserHandler_login(t *testing.T) {
	app := fiber.New()
	h := UserHandler{s: &mockUserService{}, r: app}
	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "valid request",
			body:           `{"username": "test", "password": "test"}`,
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "invalid request",
			body:           "{",
			expectedStatus: fiber.StatusUnprocessableEntity,
		}, {
			name:           "mismatched data",
			body:           `{"bad": "data"}`,
			expectedStatus: fiber.StatusUnprocessableEntity,
		},
		{
			name:           "missing username",
			body:           `{"username": "", "password": "test"}`,
			expectedStatus: fiber.StatusUnprocessableEntity,
		},
		{
			name:           "missing password",
			body:           `{"username": "test", "password": ""}`,
			expectedStatus: fiber.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Request().SetBodyString(tt.body)

			//set the request to application/json
			c.Request().Header.SetMethod("POST")
			c.Request().Header.SetContentType("application/json")

			_ = h.login(c)

			if c.Response().StatusCode() != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, c.Response().StatusCode())
			}
		})
	}
}
