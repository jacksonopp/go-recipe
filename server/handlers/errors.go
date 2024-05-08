package handlers

import "github.com/gofiber/fiber/v2"

type APIError struct {
	Code int `json:"code"`
	Msg  any `json:"msg"`

	err string
}

func NewAPIError(code int, msg any, err string) APIError {
	return APIError{Code: code, Msg: msg, err: err}
}

func (e *APIError) Error() string {
	return e.err
}

func UnprocessableEntity(items map[string]string) APIError {
	return NewAPIError(422, items, "unprocessable entity")
}

func InternalServerError() APIError {
	return NewAPIError(500, nil, "internal server error")
}

func SendError(c *fiber.Ctx, err APIError) error {
	return c.Status(err.Code).JSON(err)
}
