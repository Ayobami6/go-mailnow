package mailnow

import "fmt"

// Error represents the base error type for all Mailnow SDK errors
type Error struct {
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

// ValidationError represents input validation failures
type ValidationError struct {
	error *Error
}

// NewValidationError creates a new ValidationError
func NewValidationError(message string, err error) *ValidationError {
	return &ValidationError{
		error: &Error{
			Message: message,
			Err:     err,
		},
	}
}

func (e *ValidationError) Error() string {
	return e.error.Error()
}

func (e *ValidationError) Unwrap() error {
	return e.error.Unwrap()
}

// AuthError represents authentication failures
type AuthError struct {
	error *Error
}

// NewAuthError creates a new AuthError
func NewAuthError(message string, err error) *AuthError {
	return &AuthError{
		error: &Error{
			Message: message,
			Err:     err,
		},
	}
}

func (e *AuthError) Error() string {
	return e.error.Error()
}

func (e *AuthError) Unwrap() error {
	return e.error.Unwrap()
}

// RateLimitError represents rate limit exceeded errors
type RateLimitError struct {
	error *Error
}

// NewRateLimitError creates a new RateLimitError
func NewRateLimitError(message string, err error) *RateLimitError {
	return &RateLimitError{
		error: &Error{
			Message: message,
			Err:     err,
		},
	}
}

func (e *RateLimitError) Error() string {
	return e.error.Error()
}

func (e *RateLimitError) Unwrap() error {
	return e.error.Unwrap()
}

// ServerError represents server errors (5xx)
type ServerError struct {
	error *Error
}

// NewServerError creates a new ServerError
func NewServerError(message string, err error) *ServerError {
	return &ServerError{
		error: &Error{
			Message: message,
			Err:     err,
		},
	}
}

func (e *ServerError) Error() string {
	return e.error.Error()
}

func (e *ServerError) Unwrap() error {
	return e.error.Unwrap()
}

// ConnectionError represents network connection failures
type ConnectionError struct {
	error *Error
}

// NewConnectionError creates a new ConnectionError
func NewConnectionError(message string, err error) *ConnectionError {
	return &ConnectionError{
		error: &Error{
			Message: message,
			Err:     err,
		},
	}
}

func (e *ConnectionError) Error() string {
	return e.error.Error()
}

func (e *ConnectionError) Unwrap() error {
	return e.error.Unwrap()
}
