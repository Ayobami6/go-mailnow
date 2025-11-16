package tests

import (
	"errors"
	"testing"

	"github.com/Ayobami6/go-mailnow"
)

func TestValidateAPIKey(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
		errType error
	}{
		{
			name:    "empty API key",
			apiKey:  "",
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name:    "invalid prefix - no prefix",
			apiKey:  "invalid_key_12345",
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name:    "invalid prefix - wrong prefix",
			apiKey:  "api_live_12345",
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name:    "valid live API key",
			apiKey:  "mn_live_7e59df7ce4a14545b443837804ec9722",
			wantErr: false,
		},
		{
			name:    "valid test API key",
			apiKey:  "mn_test_abc123def456",
			wantErr: false,
		},
		{
			name:    "valid live API key - minimal",
			apiKey:  "mn_live_x",
			wantErr: false,
		},
		{
			name:    "valid test API key - minimal",
			apiKey:  "mn_test_y",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailnow.ValidateAPIKey(tt.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errType != nil {
				var validationErr *mailnow.ValidationError
				if !errors.As(err, &validationErr) {
					t.Errorf("validateAPIKey() error type = %T, want %T", err, tt.errType)
				}
			}
		})
	}
}

func TestValidateEmailAddress(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
		errType error
	}{
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name:    "invalid - no @ symbol",
			email:   "invalidemail.com",
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name:    "invalid - no domain",
			email:   "user@",
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name:    "invalid - no TLD",
			email:   "user@domain",
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name:    "invalid - no local part",
			email:   "@domain.com",
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name:    "invalid - spaces",
			email:   "user name@domain.com",
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name:    "valid - simple email",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "valid - with dots",
			email:   "first.last@example.com",
			wantErr: false,
		},
		{
			name:    "valid - with plus",
			email:   "user+tag@example.com",
			wantErr: false,
		},
		{
			name:    "valid - with hyphen",
			email:   "user-name@example.com",
			wantErr: false,
		},
		{
			name:    "valid - with numbers",
			email:   "user123@example456.com",
			wantErr: false,
		},
		{
			name:    "valid - subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailnow.ValidateEmailAddress(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEmailAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errType != nil {
				var validationErr *mailnow.ValidationError
				if !errors.As(err, &validationErr) {
					t.Errorf("validateEmailAddress() error type = %T, want %T", err, tt.errType)
				}
			}
		})
	}
}

func TestValidateEmailRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *mailnow.EmailRequest
		wantErr bool
		errType error
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name: "missing from address",
			req: &mailnow.EmailRequest{
				From:    "",
				To:      "recipient@example.com",
				Subject: "Test",
				HTML:    "<p>Test</p>",
			},
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name: "invalid from address",
			req: &mailnow.EmailRequest{
				From:    "invalid-email",
				To:      "recipient@example.com",
				Subject: "Test",
				HTML:    "<p>Test</p>",
			},
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name: "missing to address",
			req: &mailnow.EmailRequest{
				From:    "sender@example.com",
				To:      "",
				Subject: "Test",
				HTML:    "<p>Test</p>",
			},
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name: "invalid to address",
			req: &mailnow.EmailRequest{
				From:    "sender@example.com",
				To:      "invalid@",
				Subject: "Test",
				HTML:    "<p>Test</p>",
			},
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name: "missing subject",
			req: &mailnow.EmailRequest{
				From:    "sender@example.com",
				To:      "recipient@example.com",
				Subject: "",
				HTML:    "<p>Test</p>",
			},
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name: "missing HTML body",
			req: &mailnow.EmailRequest{
				From:    "sender@example.com",
				To:      "recipient@example.com",
				Subject: "Test",
				HTML:    "",
			},
			wantErr: true,
			errType: &mailnow.ValidationError{},
		},
		{
			name: "valid request - simple",
			req: &mailnow.EmailRequest{
				From:    "sender@example.com",
				To:      "recipient@example.com",
				Subject: "Test Email",
				HTML:    "<p>Test content</p>",
			},
			wantErr: false,
		},
		{
			name: "valid request - complex emails",
			req: &mailnow.EmailRequest{
				From:    "first.last+tag@mail.example.com",
				To:      "user123@subdomain.example.org",
				Subject: "Complex Test Email",
				HTML:    "<html><body><h1>Hello</h1><p>World</p></body></html>",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailnow.ValidateEmailRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEmailRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errType != nil {
				var validationErr *mailnow.ValidationError
				if !errors.As(err, &validationErr) {
					t.Errorf("validateEmailRequest() error type = %T, want %T", err, tt.errType)
				}
			}
		})
	}
}
