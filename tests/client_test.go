package tests

import (
	"errors"
	"testing"

	"github.com/Ayobami6/go-mailnow"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		apiKey      string
		expectError bool
		errorType   interface{}
	}{
		{
			name:        "empty API key should return ValidationError",
			apiKey:      "",
			expectError: true,
			errorType:   &mailnow.ValidationError{},
		},
		{
			name:        "invalid API key format (no prefix) should return ValidationError",
			apiKey:      "invalid_key_12345",
			expectError: true,
			errorType:   &mailnow.ValidationError{},
		},
		{
			name:        "invalid API key format (wrong prefix) should return ValidationError",
			apiKey:      "mn_prod_12345",
			expectError: true,
			errorType:   &mailnow.ValidationError{},
		},
		{
			name:        "valid mn_live_* API key should succeed",
			apiKey:      "mn_live_7e59df7ce4a14545b443837804ec9722",
			expectError: false,
			errorType:   nil,
		},
		{
			name:        "valid mn_test_* API key should succeed",
			apiKey:      "mn_test_7e59df7ce4a14545b443837804ec9722",
			expectError: false,
			errorType:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := mailnow.NewClient(tt.apiKey)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}

				// Check if error is of the expected type
				if !errors.As(err, &tt.errorType) {
					t.Errorf("expected error type %T, got %T", tt.errorType, err)
				}

				// Client should be nil when error occurs
				if client != nil {
					t.Errorf("expected nil client when error occurs, got %v", client)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
					return
				}

				// Client should not be nil on success
				if client == nil {
					t.Errorf("expected non-nil client on success")
				}
			}
		})
	}
}
