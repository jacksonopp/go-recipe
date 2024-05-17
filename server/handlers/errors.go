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

func BadRequest(msg string) APIError {
	return NewAPIError(400, msg)
}

func NotFound(msg map[string]string) APIError {
	return NewAPIError(404, msg)
}

func UnprocessableEntity(items map[string]string) APIError {
	return NewAPIError(422, items)
}

func Unauthorized() APIError {
	return NewAPIError(401, "Unauthorized")
}

func InternalServerError() APIError {
	return NewAPIError(500, nil)
}

func Conflict(items map[string]string) APIError {
	return NewAPIError(409, items)
}

func SendError(c *fiber.Ctx, err APIError) error {
	return c.Status(err.Code).JSON(err)
}
