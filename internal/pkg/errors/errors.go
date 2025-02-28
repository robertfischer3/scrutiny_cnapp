package errors

import (
    "fmt"
)

// ErrorType is the type of an error
type ErrorType uint

const (
    // ErrorTypeUnknown is the default error type
    ErrorTypeUnknown ErrorType = iota
    // ErrorTypeValidation is returned when there's a validation error
    ErrorTypeValidation
    // ErrorTypeDatabase is returned when there's a database error
    ErrorTypeDatabase
    // ErrorTypeNotFound is returned when a resource is not found
    ErrorTypeNotFound
)

// Error defines a standard application error
type Error struct {
    Type    ErrorType
    Message string
    Err     error
}

// Error returns the string representation of the error
func (e *Error) Error() string {
    if e.Err != nil {
        return e.Message + ": " + e.Err.Error()
    }
    return e.Message
}

// New creates a new Error
func New(errorType ErrorType, message string, err error) *Error {
    return &Error{
        Type:    errorType,
        Message: message,
        Err:     err,
    }
}

// NewValidationError creates a new validation error
func NewValidationError(message string, err error) *Error {
    return New(ErrorTypeValidation, message, err)
}

// NewDatabaseError creates a new database error
func NewDatabaseError(message string, err error) *Error {
    return New(ErrorTypeDatabase, message, err)
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, err error) *Error {
    return New(ErrorTypeNotFound, message, err)
}