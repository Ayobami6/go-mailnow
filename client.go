// Package mailnow provides a Go client library for the Mailnow email API.
//
// The Mailnow SDK enables customers to send emails programmatically through
// the Mailnow SaaS service (https://api.mailnow.xyz) using API key authentication.
//
// Basic usage:
//
//	client, err := mailnow.NewClient("mn_live_your_api_key")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	ctx := context.Background()
//	req := &mailnow.EmailRequest{
//	    From:    "sender@example.com",
//	    To:      "recipient@example.com",
//	    Subject: "Hello",
//	    HTML:    "<h1>Hello World</h1>",
//	}
//
//	resp, err := client.SendEmail(ctx, req)
//	if err != nil {
//	    // Handle error (ValidationError, AuthError, etc.)
//	    log.Fatal(err)
//	}
//	fmt.Printf("Email sent: %s\n", resp.MessageID)
package mailnow

import (
	"context"
	"encoding/json"
	"net/http"
)

// Client represents a Mailnow API client for sending emails.
//
// The Client handles authentication, request validation, and communication
// with the Mailnow API. It maintains an HTTP client with connection pooling
// for efficient request handling.
//
// A Client should be created using NewClient and can be safely reused
// across multiple goroutines for sending multiple emails.
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewClient creates and initializes a new Mailnow API client.
//
// The apiKey parameter must be a valid Mailnow API key starting with
// either "mn_live_" (for production) or "mn_test_" (for testing).
//
// Returns a configured Client ready to send emails, or an error if
// the API key is invalid.
//
// Example:
//
//	client, err := mailnow.NewClient("mn_live_7e59df7ce4a14545b443837804ec9722")
//	
	// Validate API key
	if err := ValidateAPIKey(apiKey); err != nil {
		return nil, err
	}

	// Initialize HTTP client with timeout configuration
	httpClient := &http.Client{
		Timeout: RequestTimeout,
	}

	// Create and return the client
	return &Client{
		apiKey:     apiKey,
		httpClient: httpClient,
		baseURL:    APIBaseURL,
	}, nil
}

// SendEmail sends an email via the Mailnow API
func (c *Client) SendEmail(ctx context.Context, req *EmailRequest) (*EmailResponse, error) {
	// Validate email request
	if err := ValidateEmailRequest(req); err != nil {
		return nil, err
	}

	// Build full URL
	url := c.baseURL + EmailSendEndpoint

	// Make HTTP POST request
	resp, err := MakeRequest(ctx, c.httpClient, "POST", url, c.apiKey, req)
	if err != nil {
		return nil, err
	}

	// Handle response
	body, err := HandleResponse(resp)
	if err != nil {
		return nil, err
	}

	// Parse successful response JSON into EmailResponse struct
	var emailResp EmailResponse
	if err := json.Unmarshal(body, &emailResp); err != nil {
		return nil, NewServerError("failed to parse response", err)
	}

	return &emailResp, nil
}
