package services

import "errors"

var (
	ErrTimeout          = errors.New("timeout")
	ErrTimeoutNoMessage = errors.New("timeout with no message")
	ErrCommit           = errors.New("commit error")

	ErrUnknown = errors.New("unknown error")
	//ErrUserAlreadyExists = errors.New("user already exists")
	//ErrUserNotFound = errors.New("user not found")
	//ErrInvalidPassword = errors.New("invalid password")
	//
	ErrRecipeNotFound      = errors.New("recipe not found")
	ErrIngredientNotFound  = errors.New("ingredient not found")
	ErrIngredientConflict  = errors.New("ingredient conflict")
	ErrInstructionNotFound = errors.New("instruction not found")
	ErrInstructionConflict = errors.New("instruction conflict")
	//
	//ErrTagConflict = errors.New("tag conflict")
	//
	//ErrSessionNotFound = errors.New("session not found")
	//ErrSessionExpired = errors.New("session expired")

	ErrTagNotFound = errors.New("tag not found")
	ErrTagConflict = errors.New("tag conflict")
)
