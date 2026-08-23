package domain

import "errors"

var (
	ErrInvalid           = errors.New("invalid domain input")
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("conflict")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrForbidden         = errors.New("forbidden")
)
