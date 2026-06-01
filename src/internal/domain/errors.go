package domain

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("already exists")
	ErrInvalidInput = errors.New("invalid input")
	ErrBusinessRule = errors.New("business rule violation")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type DomainError struct {
	Code    error
	Message string
}

func (e *DomainError) Error() string { return e.Message }
func (e *DomainError) Unwrap() error { return e.Code }

func NewNotFound(msg string) *DomainError {
	return &DomainError{Code: ErrNotFound, Message: msg}
}
func NewConflict(msg string) *DomainError {
	return &DomainError{Code: ErrConflict, Message: msg}
}
func NewInvalidInput(msg string) *DomainError {
	return &DomainError{Code: ErrInvalidInput, Message: msg}
}
func NewBusinessRule(msg string) *DomainError {
	return &DomainError{Code: ErrBusinessRule, Message: msg}
}
