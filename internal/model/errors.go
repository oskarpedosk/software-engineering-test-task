package model

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidEmail = errors.New("invalid email format")
	ErrEmptyField   = errors.New("required field is empty")
)
