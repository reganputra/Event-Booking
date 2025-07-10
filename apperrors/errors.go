package apperrors

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrBadRequest     = errors.New("bad request")
	ErrForbidden      = errors.New("forbidden")
	ErrConflict       = errors.New("conflict")
	ErrInternalServer = errors.New("internal server error")
	ErrInvalidInput   = errors.New("invalid input")
	ErrAlreadyExists  = errors.New("already exists")
)
