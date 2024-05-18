package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type APIError struct {
	Code int `json:"code"`
	Msg  any `json:"msg"`
}

func NewAPIError(code int, msg any) APIError {
	return APIError{Code: code, Msg: msg}
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %v", e.Msg)
}

// BadRequest returns a 400 Bad Request error with the given message.
func BadRequest(msg string) APIError {
	return NewAPIError(400, msg)
}

// Unauthorized returns a 401 Unauthorized error.
func Unauthorized() APIError {
	return NewAPIError(401, "Unauthorized")
}

// NotFound returns a 404 Not Found error with the given message.
func NotFound(msg map[string]string) APIError {
	return NewAPIError(404, msg)
}

// Conflict returns a 409 Conflict error with the given items.
func Conflict(items map[string]string) APIError {
	return NewAPIError(409, items)
}

// UnprocessableEntity returns a 422 Unprocessable Entity error with the given items.
func UnprocessableEntity(items map[string]string) APIError {
	return NewAPIError(422, items)
}

// InternalServerError returns a 500 Internal Server Error.
func InternalServerError() APIError {
	return NewAPIError(500, nil)
}

// SendError sends an APIError as a JSON response.
func SendError(c *fiber.Ctx, err APIError) error {
	return c.Status(err.Code).JSON(err)
}
