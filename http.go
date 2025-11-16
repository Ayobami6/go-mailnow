package mailnow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MakeRequest builds and sends an HTTP request with proper headers
func MakeRequest(ctx context.Context, client *http.Client, method, url, apiKey string, body interface{}) (*http.Response, error) {
	// Encode request body as JSON
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, NewValidationError("failed to encode request body", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, NewConnectionError("failed to create request", err)
	}

	// Add required headers
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewConnectionError("failed to send request", err)
	}

	return resp, nil
}

// HandleResponse processes HTTP responses and maps status codes to error types
func HandleResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewConnectionError("failed to read response body", err)
	}

	// Handle successful responses (2xx)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, nil
	}

	// Parse error response
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		// If we can't parse the error response, create a generic error message
		return nil, mapStatusCodeToError(resp.StatusCode, string(body))
	}

	// Map status code to appropriate error type with parsed message
	errorMessage := errResp.Error.Message
	if errorMessage == "" {
		errorMessage = fmt.Sprintf("API request failed with status %d", resp.StatusCode)
	}

	return nil, mapStatusCodeToError(resp.StatusCode, errorMessage)
}

// mapStatusCodeToError maps HTTP status codes to specific error types
func mapStatusCodeToError(statusCode int, message string) error {
	switch statusCode {
	case 400:
		return NewValidationError(message, nil)
	case 401:
		return NewAuthError(message, nil)
	case 429:
		return NewRateLimitError(message, nil)
	default:
		if statusCode >= 500 {
			return NewServerError(message, nil)
		}
		return NewServerError(fmt.Sprintf("unexpected status code %d: %s", statusCode, message), nil)
	}
}
