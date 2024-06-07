package services

import "errors"

var (
	// General errors

	// ErrUnknown is returned when an error is not recognized
	ErrUnknown = errors.New("unknown error")

	// ErrTimeout is returned when an operation times out
	ErrTimeout = errors.New("timeout")

	// ErrTimeoutNoMessage is returned when an operation times out with no message
	ErrTimeoutNoMessage = errors.New("timeout with no message")

	// ErrCommit is returned when a transaction commit fails
	ErrCommit = errors.New("commit error")

	// ErrUnauthorized is returned when a user is not authorized
	ErrUnauthorized = errors.New("unauthorized")

	// Auth errors

	// ErrUserAlreadyExists is returned when a user already exists
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")

	// ErrPasswordMismatch is returned when a password does not match
	ErrPasswordMismatch = errors.New("password mismatch")

	// ErrInvalidPassword is returned when a password is invalid
	ErrInvalidPassword = errors.New("invalid password")

	// Recipe errors

	// ErrRecipeNotFound is returned when a recipe is not found
	ErrRecipeNotFound = errors.New("recipe not found")

	// ErrIngredientNotFound is returned when an ingredient is not found
	ErrIngredientNotFound = errors.New("ingredient not found")

	// ErrIngredientConflict is returned when an ingredient conflict occurs
	ErrIngredientConflict = errors.New("ingredient conflict")

	// ErrInstructionNotFound is returned when an instruction is not found
	ErrInstructionNotFound = errors.New("instruction not found")

	// ErrInstructionConflict is returned when an instruction conflict occurs
	ErrInstructionConflict = errors.New("instruction conflict")

	// Tag Errors

	// ErrTagNotFound is returned when a tag is not found
	ErrTagNotFound = errors.New("tag not found")

	// ErrTagConflict is returned when a tag conflict occurs
	ErrTagConflict = errors.New("tag conflict")

	// Session errors

	// ErrSessionNotFound is returned when a session is not found
	ErrSessionNotFound = errors.New("session not found")

	// ErrSessionExpired is returned when a session has expired
	ErrSessionExpired = errors.New("session expired")
)
