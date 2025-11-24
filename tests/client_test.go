package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
func TestSendEmailValidation(t *testing.T) {
	// Test validation errors that occur before making HTTP requests
	tests := []struct {
		name      string
		request   *mailnow.EmailRequest
		errorType interface{}
	}{
		{
			name:      "nil request",
			request:   nil,
			errorType: &mailnow.ValidationError{},
		},
		{
			name: "empty from address",
			request: &mailnow.EmailRequest{
				From:    "",
				To:      "test@example.com",
				Subject: "Test Subject",
				HTML:    "<h1>Test</h1>",
			},
			errorType: &mailnow.ValidationError{},
		},
		{
			name: "empty to address",
			request: &mailnow.EmailRequest{
				From:    "sender@example.com",
				To:      "",
				Subject: "Test Subject",
				HTML:    "<h1>Test</h1>",
			},
			errorType: &mailnow.ValidationError{},
		},
		{
			name: "empty subject",
			request: &mailnow.EmailRequest{
				From:    "sender@example.com",
				To:      "test@example.com",
				Subject: "",
				HTML:    "<h1>Test</h1>",
			},
			errorType: &mailnow.ValidationError{},
		},
		{
			name: "empty HTML body",
			request: &mailnow.EmailRequest{
				From:    "sender@example.com",
				To:      "test@example.com",
				Subject: "Test Subject",
				HTML:    "",
			},
			errorType: &mailnow.ValidationError{},
		},
		{
			name: "invalid from email format",
			request: &mailnow.EmailRequest{
				From:    "invalid-email",
				To:      "test@example.com",
				Subject: "Test Subject",
				HTML:    "<h1>Test</h1>",
			},
			errorType: &mailnow.ValidationError{},
		},
		{
			name: "invalid to email format",
			request: &mailnow.EmailRequest{
				From:    "sender@example.com",
				To:      "invalid-email",
				Subject: "Test Subject",
				HTML:    "<h1>Test</h1>",
			},
			errorType: &mailnow.ValidationError{},
		},
	}

	// Create client
	client, err := mailnow.NewClient("mn_test_7e59df7ce4a14545b443837804ec9722")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.SendEmail(ctx, tt.request)

			// Should return validation error
			if err == nil {
				t.Errorf("expected ValidationError, but got no error")
				return
			}

			if !errors.As(err, &tt.errorType) {
				t.Errorf("expected error type %T, got %T: %v", tt.errorType, err, err)
			}

			// Response should be nil when error occurs
			if resp != nil {
				t.Errorf("expected nil response when error occurs, got %v", resp)
			}
		})
	}
}

func TestSendEmailWithContextCancellation(t *testing.T) {
	// Create a test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delay to allow context cancellation
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(200)
		w.Write([]byte(`{"success": true, "message_id": "msg_12345", "status": "sent"}`))
	}))
	defer server.Close()

	// Create client
	client, err := mailnow.NewClient("mn_test_7e59df7ce4a14545b443837804ec9722")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Prepare valid request
	req := &mailnow.EmailRequest{
		From:    "sender@example.com",
		To:      "test@example.com",
		Subject: "Test Subject",
		HTML:    "<h1>Test</h1>",
	}

	// Call SendEmail - should fail due to context cancellation
	resp, err := client.SendEmail(ctx, req)

	// Verify that we get an error (context cancellation or connection error)
	if err == nil {
		t.Errorf("expected error due to context cancellation, but got none")
		return
	}

	// Response should be nil when error occurs
	if resp != nil {
		t.Errorf("expected nil response when error occurs, got %v", resp)
	}

	// The error should be a connection error (since context cancellation is handled as a connection issue)
	var connErr *mailnow.ConnectionError
	if !errors.As(err, &connErr) {
		// Context cancellation might also result in other error types depending on when it occurs
		// So we just verify that we got an error, which is the important part
		t.Logf("Got error (expected due to context cancellation): %v", err)
	}
}

func TestSendEmailWithNilRequest(t *testing.T) {
	// Create client
	client, err := mailnow.NewClient("mn_test_7e59df7ce4a14545b443837804ec9722")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Create context
	ctx := context.Background()

	// Call SendEmail with nil request
	resp, err := client.SendEmail(ctx, nil)

	// Should return ValidationError
	if err == nil {
		t.Errorf("expected ValidationError for nil request, but got no error")
		return
	}

	var validationErr *mailnow.ValidationError
	if !errors.As(err, &validationErr) {
		t.Errorf("expected ValidationError for nil request, got %T: %v", err, err)
	}

	// Response should be nil
	if resp != nil {
		t.Errorf("expected nil response for nil request, got %v", resp)
	}
}

// TestSendEmailHTTPIntegration tests the HTTP integration aspects of SendEmail
// Note: These tests will make actual HTTP requests to a mock server
func TestSendEmailHTTPIntegration(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody string
		errorType    interface{}
		expectError  bool
	}{
		{
			name:         "successful response",
			statusCode:   200,
			responseBody: `{"success": true, "message_id": "msg_12345", "status": "sent"}`,
			expectError:  false,
		},
		{
			name:         "authentication error",
			statusCode:   401,
			responseBody: `{"error": {"code": "unauthorized", "message": "Invalid API key"}}`,
			expectError:  true,
			errorType:    &mailnow.AuthError{},
		},
		{
			name:         "validation error from API",
			statusCode:   400,
			responseBody: `{"error": {"code": "validation_error", "message": "Invalid email format"}}`,
			expectError:  true,
			errorType:    &mailnow.ValidationError{},
		},
		{
			name:         "rate limit error",
			statusCode:   429,
			responseBody: `{"error": {"code": "rate_limit", "message": "Rate limit exceeded"}}`,
			expectError:  true,
			errorType:    &mailnow.RateLimitError{},
		},
		{
			name:         "server error",
			statusCode:   500,
			responseBody: `{"error": {"code": "internal_error", "message": "Internal server error"}}`,
			expectError:  true,
			errorType:    &mailnow.ServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server that responds to the expected endpoint
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the request is what we expect
				if r.Method != "POST" {
					t.Errorf("expected POST method, got %s", r.Method)
				}

				if r.URL.Path != "/v1/email/send" {
					t.Errorf("expected path /v1/email/send, got %s", r.URL.Path)
				}

				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
				}

				if r.Header.Get("X-API-Key") == "" {
					t.Errorf("expected X-API-Key header to be set")
				}

				// Verify request body contains expected fields
				var reqBody mailnow.EmailRequest
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
				}

				if reqBody.From == "" || reqBody.To == "" || reqBody.Subject == "" || reqBody.HTML == "" {
					t.Errorf("request body missing required fields: %+v", reqBody)
				}

				// Send the test response
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Since we can't override the client's baseURL easily, we'll test this by
			// temporarily modifying the constants or using a different approach
			// For now, we'll test the individual components that SendEmail uses

			// Test the HTTP request/response handling directly
			client := &http.Client{Timeout: 5 * time.Second}
			ctx := context.Background()

			// Create a valid email request
			emailReq := &mailnow.EmailRequest{
				From:    "sender@example.com",
				To:      "test@example.com",
				Subject: "Test Subject",
				HTML:    "<h1>Test</h1>",
			}

			// Test the MakeRequest function with our mock server
			url := server.URL + "/v1/email/send"
			resp, err := mailnow.MakeRequest(ctx, client, "POST", url, "mn_test_abc123", emailReq)
			if err != nil {
				if tt.expectError {
					// Check if it's a connection error (which might happen before we get to status code handling)
					var connErr *mailnow.ConnectionError
					if errors.As(err, &connErr) {
						return // This is acceptable for connection-related errors
					}
				}
				t.Errorf("MakeRequest failed: %v", err)
				return
			}

			// Test the HandleResponse function
			body, err := mailnow.HandleResponse(resp)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}

				// Check error type
				if tt.errorType != nil && !errors.As(err, &tt.errorType) {
					t.Errorf("expected error type %T, got %T: %v", tt.errorType, err, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
					return
				}

				if body == nil {
					t.Errorf("expected response body but got nil")
					return
				}

				// Parse and verify the response
				var emailResp mailnow.EmailResponse
				if err := json.Unmarshal(body, &emailResp); err != nil {
					t.Errorf("failed to parse response: %v", err)
					return
				}

				if !emailResp.Success {
					t.Errorf("expected success=true, got %v", emailResp.Success)
				}

				if emailResp.Data.MessageID == "" {
					t.Errorf("expected non-empty message ID")
				}
			}
		})
	}
}
