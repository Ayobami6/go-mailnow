package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Ayobami6/go-mailnow"
)

// TestErrorCreationAndFormatting tests error creation and message formatting
func TestErrorCreationAndFormatting(t *testing.T) {
	tests := []struct {
		name            string
		createError     func() error
		expectedMessage string
	}{
		{
			name: "ValidationError without wrapped error",
			createError: func() error {
				return mailnow.NewValidationError("invalid email format", nil)
			},
			expectedMessage: "invalid email format",
		},
		{
			name: "ValidationError with wrapped error",
			createError: func() error {
				baseErr := fmt.Errorf("field is empty")
				return mailnow.NewValidationError("validation failed", baseErr)
			},
			expectedMessage: "validation failed: field is empty",
		},
		{
			name: "AuthError without wrapped error",
			createError: func() error {
				return mailnow.NewAuthError("invalid API key", nil)
			},
			expectedMessage: "invalid API key",
		},
		{
			name: "AuthError with wrapped error",
			createError: func() error {
				baseErr := fmt.Errorf("unauthorized")
				return mailnow.NewAuthError("authentication failed", baseErr)
			},
			expectedMessage: "authentication failed: unauthorized",
		},
		{
			name: "RateLimitError without wrapped error",
			createError: func() error {
				return mailnow.NewRateLimitError("rate limit exceeded", nil)
			},
			expectedMessage: "rate limit exceeded",
		},
		{
			name: "RateLimitError with wrapped error",
			createError: func() error {
				baseErr := fmt.Errorf("too many requests")
				return mailnow.NewRateLimitError("rate limit hit", baseErr)
			},
			expectedMessage: "rate limit hit: too many requests",
		},
		{
			name: "ServerError without wrapped error",
			createError: func() error {
				return mailnow.NewServerError("internal server error", nil)
			},
			expectedMessage: "internal server error",
		},
		{
			name: "ServerError with wrapped error",
			createError: func() error {
				baseErr := fmt.Errorf("database connection failed")
				return mailnow.NewServerError("server error occurred", baseErr)
			},
			expectedMessage: "server error occurred: database connection failed",
		},
		{
			name: "ConnectionError without wrapped error",
			createError: func() error {
				return mailnow.NewConnectionError("connection timeout", nil)
			},
			expectedMessage: "connection timeout",
		},
		{
			name: "ConnectionError with wrapped error",
			createError: func() error {
				baseErr := fmt.Errorf("network unreachable")
				return mailnow.NewConnectionError("connection failed", baseErr)
			},
			expectedMessage: "connection failed: network unreachable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.createError()
			if err.Error() != tt.expectedMessage {
				t.Errorf("expected message %q, got %q", tt.expectedMessage, err.Error())
			}
		})
	}
}

// TestErrorUnwrapping tests error unwrapping with errors.Unwrap()
func TestErrorUnwrapping(t *testing.T) {
	baseErr := fmt.Errorf("base error")

	tests := []struct {
		name         string
		createError  func() error
		shouldUnwrap bool
	}{
		{
			name: "ValidationError with wrapped error",
			createError: func() error {
				return mailnow.NewValidationError("validation failed", baseErr)
			},
			shouldUnwrap: true,
		},
		{
			name: "ValidationError without wrapped error",
			createError: func() error {
				return mailnow.NewValidationError("validation failed", nil)
			},
			shouldUnwrap: false,
		},
		{
			name: "AuthError with wrapped error",
			createError: func() error {
				return mailnow.NewAuthError("auth failed", baseErr)
			},
			shouldUnwrap: true,
		},
		{
			name: "AuthError without wrapped error",
			createError: func() error {
				return mailnow.NewAuthError("auth failed", nil)
			},
			shouldUnwrap: false,
		},
		{
			name: "RateLimitError with wrapped error",
			createError: func() error {
				return mailnow.NewRateLimitError("rate limit", baseErr)
			},
			shouldUnwrap: true,
		},
		{
			name: "RateLimitError without wrapped error",
			createError: func() error {
				return mailnow.NewRateLimitError("rate limit", nil)
			},
			shouldUnwrap: false,
		},
		{
			name: "ServerError with wrapped error",
			createError: func() error {
				return mailnow.NewServerError("server error", baseErr)
			},
			shouldUnwrap: true,
		},
		{
			name: "ServerError without wrapped error",
			createError: func() error {
				return mailnow.NewServerError("server error", nil)
			},
			shouldUnwrap: false,
		},
		{
			name: "ConnectionError with wrapped error",
			createError: func() error {
				return mailnow.NewConnectionError("connection error", baseErr)
			},
			shouldUnwrap: true,
		},
		{
			name: "ConnectionError without wrapped error",
			createError: func() error {
				return mailnow.NewConnectionError("connection error", nil)
			},
			shouldUnwrap: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.createError()
			unwrapped := errors.Unwrap(err)

			if tt.shouldUnwrap {
				if unwrapped == nil {
					t.Error("expected error to unwrap, but got nil")
				}
				if unwrapped != baseErr {
					t.Errorf("expected unwrapped error to be baseErr, got %v", unwrapped)
				}
			} else {
				if unwrapped != nil {
					t.Errorf("expected nil when unwrapping, got %v", unwrapped)
				}
			}
		})
	}
}

// TestErrorTypeAssertions tests type assertions with errors.As()
func TestErrorTypeAssertions(t *testing.T) {
	tests := []struct {
		name        string
		createError func() error
		assertType  string
	}{
		{
			name: "ValidationError type assertion",
			createError: func() error {
				return mailnow.NewValidationError("validation failed", nil)
			},
			assertType: "ValidationError",
		},
		{
			name: "AuthError type assertion",
			createError: func() error {
				return mailnow.NewAuthError("auth failed", nil)
			},
			assertType: "AuthError",
		},
		{
			name: "RateLimitError type assertion",
			createError: func() error {
				return mailnow.NewRateLimitError("rate limit", nil)
			},
			assertType: "RateLimitError",
		},
		{
			name: "ServerError type assertion",
			createError: func() error {
				return mailnow.NewServerError("server error", nil)
			},
			assertType: "ServerError",
		},
		{
			name: "ConnectionError type assertion",
			createError: func() error {
				return mailnow.NewConnectionError("connection error", nil)
			},
			assertType: "ConnectionError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.createError()

			switch tt.assertType {
			case "ValidationError":
				var validationErr *mailnow.ValidationError
				if !errors.As(err, &validationErr) {
					t.Error("expected error to be ValidationError")
				}
				// Verify we can't assert to other types
				var authErr *mailnow.AuthError
				if errors.As(err, &authErr) {
					t.Error("ValidationError should not assert to AuthError")
				}

			case "AuthError":
				var authErr *mailnow.AuthError
				if !errors.As(err, &authErr) {
					t.Error("expected error to be AuthError")
				}
				// Verify we can't assert to other types
				var validationErr *mailnow.ValidationError
				if errors.As(err, &validationErr) {
					t.Error("AuthError should not assert to ValidationError")
				}

			case "RateLimitError":
				var rateLimitErr *mailnow.RateLimitError
				if !errors.As(err, &rateLimitErr) {
					t.Error("expected error to be RateLimitError")
				}
				// Verify we can't assert to other types
				var validationErr *mailnow.ValidationError
				if errors.As(err, &validationErr) {
					t.Error("RateLimitError should not assert to ValidationError")
				}

			case "ServerError":
				var serverErr *mailnow.ServerError
				if !errors.As(err, &serverErr) {
					t.Error("expected error to be ServerError")
				}
				// Verify we can't assert to other types
				var validationErr *mailnow.ValidationError
				if errors.As(err, &validationErr) {
					t.Error("ServerError should not assert to ValidationError")
				}

			case "ConnectionError":
				var connErr *mailnow.ConnectionError
				if !errors.As(err, &connErr) {
					t.Error("expected error to be ConnectionError")
				}
				// Verify we can't assert to other types
				var validationErr *mailnow.ValidationError
				if errors.As(err, &validationErr) {
					t.Error("ConnectionError should not assert to ValidationError")
				}
			}
		})
	}
}
