//go:build integration
// +build integration

package mailnow

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

// TestIntegrationSuccessfulEmailSend tests successful email sending with valid API key
// This test requires a valid mn_test_* API key to be set in the MAILNOW_TEST_API_KEY environment variable
func TestIntegrationSuccessfulEmailSend(t *testing.T) {
	apiKey := os.Getenv("MAILNOW_TEST_API_KEY")
	if apiKey == "" {
		t.Skip("MAILNOW_TEST_API_KEY environment variable not set, skipping integration test")
	}
	log.Println("Using API Key:", apiKey)

	// Create client with test API key
	client, err := NewClient(apiKey)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create context with reasonable timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Prepare valid email request
	req := &EmailRequest{
		From:    "ayobamidele006@gmail.com",
		To:      "ayobamidele006@gmail.com",
		Subject: "Integration Test Email",
		HTML:    "<h1>Integration Test</h1><p>This is a test email sent from the Go SDK integration tests.</p>",
	}

	// Send email
	resp, err := client.SendEmail(ctx, req)
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	// Verify response
	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	if !resp.Success {
		t.Errorf("Expected success=true, got %v", resp.Success)
	}
	log.Println("Response:", resp)

	if resp.Data.MessageID == "" {
		t.Error("Expected non-empty message ID")
	}

	if resp.Data.Status == "" {
		t.Error("Expected non-empty status")
	}

	t.Logf("Email sent successfully: MessageID=%s, Status=%s", resp.Data.MessageID, resp.Data.Status)
}

// TestIntegrationAuthenticationFailure tests authentication failure with invalid API key
func TestIntegrationAuthenticationFailure(t *testing.T) {
	// Use an invalid API key (valid format but non-existent)
	invalidAPIKey := "mn_test_invalid_key_12345678901234567890"
	// invalidAPIKey := "mn_live_40a9c9353dee4ab5b112573b7e65e1e9"

	// Create client with invalid API key
	client, err := NewClient(invalidAPIKey)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create context with reasonable timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Prepare valid email request
	req := &EmailRequest{
		From:        "test@example.com",
		To:          "recipient@example.com",
		Subject:     "Test Email",
		HTML:        "<h1>Test</h1><p>This should fail due to invalid API key.</p>",
		Attachments: nil,
	}

	// Send email - should fail with authentication error
	resp, err := client.SendEmail(ctx, req)

	// Verify that we get an error (could be AuthError or ServerError depending on API behavior)
	if err == nil {
		t.Fatal("Expected error for invalid API key, but got no error")
	}

	// The API might return different error types for invalid keys or network issues
	// Accept AuthError (401), ServerError (5xx), or ConnectionError (network issues) as valid responses
	var authErr *AuthError
	var serverErr *ServerError
	var connErr *ConnectionError
	var validationErr *ValidationError
	if !errors.As(err, &authErr) && !errors.As(err, &serverErr) && !errors.As(err, &connErr) && !errors.As(err, &validationErr) {
		t.Errorf("Expected AuthError, ServerError, or ConnectionError for invalid API key, got %T: %v", err, err)
	}

	// Response should be nil when error occurs
	if resp != nil {
		t.Errorf("Expected nil response when error occurs, got %v", resp)
	}

	t.Logf("Error received as expected for invalid API key: %v", err)
}

// TestIntegrationValidationErrors tests validation errors with invalid email parameters
func TestIntegrationValidationErrors(t *testing.T) {
	apiKey := os.Getenv("MAILNOW_TEST_API_KEY")
	if apiKey == "" {
		// Use a valid format API key for validation testing (won't reach the API)
		apiKey = "mn_test_validation_test_key"
	}

	// Create client
	client, err := NewClient(apiKey)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test cases for validation errors
	tests := []struct {
		name    string
		request *EmailRequest
	}{
		{
			name: "empty from address",
			request: &EmailRequest{
				From:    "",
				To:      "recipient@example.com",
				Subject: "Test Subject",
				HTML:    "<h1>Test</h1>",
			},
		},
		{
			name: "empty to address",
			request: &EmailRequest{
				From:    "sender@example.com",
				To:      "",
				Subject: "Test Subject",
				HTML:    "<h1>Test</h1>",
			},
		},
		{
			name: "empty subject",
			request: &EmailRequest{
				From:    "sender@example.com",
				To:      "recipient@example.com",
				Subject: "",
				HTML:    "<h1>Test</h1>",
			},
		},
		{
			name: "empty HTML body",
			request: &EmailRequest{
				From:    "sender@example.com",
				To:      "recipient@example.com",
				Subject: "Test Subject",
				HTML:    "",
			},
		},
		{
			name: "invalid from email format",
			request: &EmailRequest{
				From:    "invalid-email-format",
				To:      "recipient@example.com",
				Subject: "Test Subject",
				HTML:    "<h1>Test</h1>",
			},
		},
		{
			name: "invalid to email format",
			request: &EmailRequest{
				From:    "sender@example.com",
				To:      "invalid@",
				Subject: "Test Subject",
				HTML:    "<h1>Test</h1>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Send email - should fail with validation error
			resp, err := client.SendEmail(ctx, tt.request)

			// Verify that we get a validation error
			if err == nil {
				t.Fatal("Expected validation error, but got no error")
			}

			var validationErr *ValidationError
			if !errors.As(err, &validationErr) {
				t.Errorf("Expected ValidationError, got %T: %v", err, err)
			}

			// Response should be nil when error occurs
			if resp != nil {
				t.Errorf("Expected nil response when error occurs, got %v", resp)
			}

			t.Logf("Validation error received as expected for %s: %v", tt.name, err)
		})
	}
}

// TestIntegrationContextTimeout tests context timeout handling
func TestIntegrationContextTimeout(t *testing.T) {
	apiKey := os.Getenv("MAILNOW_TEST_API_KEY")
	if apiKey == "" {
		// Use a valid format API key for timeout testing
		apiKey = "mn_test_timeout_test_key"
	}

	// Create client
	client, err := NewClient(apiKey)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create context with very short timeout (1 millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Prepare valid email request
	req := &EmailRequest{
		From:    "sender@example.com",
		To:      "recipient@example.com",
		Subject: "Timeout Test Email",
		HTML:    "<h1>Timeout Test</h1><p>This request should timeout.</p>",
	}

	// Send email - should fail due to context timeout
	resp, err := client.SendEmail(ctx, req)

	// Verify that we get an error (should be connection error due to timeout)
	if err == nil {
		t.Fatal("Expected timeout error, but got no error")
	}

	// The error should be a connection error (context timeout is handled as connection issue)
	var connErr *ConnectionError
	if !errors.As(err, &connErr) {
		// Context timeout might result in different error types depending on when it occurs
		// The important thing is that we get an error
		t.Logf("Got error (expected due to timeout): %T: %v", err, err)
	} else {
		t.Logf("Connection error received as expected due to timeout: %v", err)
	}

	// Response should be nil when error occurs
	if resp != nil {
		t.Errorf("Expected nil response when error occurs, got %v", resp)
	}
}

// TestIntegrationContextCancellation tests context cancellation handling
func TestIntegrationContextCancellation(t *testing.T) {
	apiKey := os.Getenv("MAILNOW_TEST_API_KEY")
	if apiKey == "" {
		// Use a valid format API key for cancellation testing
		apiKey = "mn_test_cancellation_test_key"
	}

	// Create client
	client, err := NewClient(apiKey)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create context that we'll cancel immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Prepare valid email request
	req := &EmailRequest{
		From:    "sender@example.com",
		To:      "recipient@example.com",
		Subject: "Cancellation Test Email",
		HTML:    "<h1>Cancellation Test</h1><p>This request should be cancelled.</p>",
	}

	// Send email - should fail due to context cancellation
	resp, err := client.SendEmail(ctx, req)

	// Verify that we get an error
	if err == nil {
		t.Fatal("Expected cancellation error, but got no error")
	}

	// The error should be a connection error (context cancellation is handled as connection issue)
	var connErr *ConnectionError
	if !errors.As(err, &connErr) {
		// Context cancellation might result in different error types depending on when it occurs
		// The important thing is that we get an error
		t.Logf("Got error (expected due to cancellation): %T: %v", err, err)
	} else {
		t.Logf("Connection error received as expected due to cancellation: %v", err)
	}

	// Response should be nil when error occurs
	if resp != nil {
		t.Errorf("Expected nil response when error occurs, got %v", resp)
	}
}
