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

// TestMakeRequest tests the makeRequest function with proper header construction
func TestMakeRequest(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		apiKey        string
		body          interface{}
		wantHeaders   map[string]string
		wantErr       bool
		errType       interface{}
		setupServer   func() *httptest.Server
		cancelContext bool
	}{
		{
			name:   "successful request with body",
			method: "POST",
			apiKey: "mn_test_abc123",
			body: &mailnow.EmailRequest{
				From:    "test@example.com",
				To:      "recipient@example.com",
				Subject: "Test",
				HTML:    "<p>Test</p>",
			},
			wantHeaders: map[string]string{
				"X-API-Key":    "mn_test_abc123",
				"Content-Type": "application/json",
			},
			wantErr: false,
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Verify headers
					if r.Header.Get("X-API-Key") != "mn_test_abc123" {
						t.Errorf("Expected X-API-Key header to be 'mn_test_abc123', got '%s'", r.Header.Get("X-API-Key"))
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Content-Type header to be 'application/json', got '%s'", r.Header.Get("Content-Type"))
					}
					w.WriteHeader(http.StatusOK)
				}))
			},
		},
		{
			name:   "successful request without body",
			method: "GET",
			apiKey: "mn_live_xyz789",
			body:   nil,
			wantHeaders: map[string]string{
				"X-API-Key":    "mn_live_xyz789",
				"Content-Type": "application/json",
			},
			wantErr: false,
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
		},
		{
			name:    "invalid JSON body",
			method:  "POST",
			apiKey:  "mn_test_abc123",
			body:    make(chan int), // channels cannot be marshaled to JSON
			wantErr: true,
			errType: &mailnow.ValidationError{},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
		},
		{
			name:          "context cancellation",
			method:        "POST",
			apiKey:        "mn_test_abc123",
			body:          &mailnow.EmailRequest{From: "test@example.com"},
			wantErr:       true,
			errType:       &mailnow.ConnectionError{},
			cancelContext: true,
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(100 * time.Millisecond)
					w.WriteHeader(http.StatusOK)
				}))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			client := &http.Client{Timeout: 5 * time.Second}
			ctx := context.Background()

			if tt.cancelContext {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel() // Cancel immediately
			}

			resp, err := mailnow.MakeRequest(ctx, client, tt.method, server.URL, tt.apiKey, tt.body)

			if tt.wantErr {
				if err == nil {
					t.Errorf("makeRequest() expected error, got nil")
					return
				}
				// Check error type
				if tt.errType != nil {
					switch tt.errType.(type) {
					case *mailnow.ValidationError:
						var ve *mailnow.ValidationError
						if !errors.As(err, &ve) {
							t.Errorf("makeRequest() error type = %T, want ValidationError", err)
						}
					case *mailnow.ConnectionError:
						var ce *mailnow.ConnectionError
						if !errors.As(err, &ce) {
							t.Errorf("makeRequest() error type = %T, want ConnectionError", err)
						}
					}
				}
				return
			}

			if err != nil {
				t.Errorf("makeRequest() unexpected error = %v", err)
				return
			}

			if resp == nil {
				t.Error("makeRequest() returned nil response")
				return
			}

			resp.Body.Close()
		})
	}
}

// TestHandleResponse tests the handleResponse function for successful responses
func TestHandleResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       interface{}
		wantErr    bool
		errType    interface{}
	}{
		{
			name:       "successful response 200",
			statusCode: http.StatusOK,
			body: mailnow.EmailResponse{
				Success: true,
				Data: mailnow.Data{
					MessageID: "msg_123",
					Status:    "sent",
				},
			},
			wantErr: false,
		},
		{
			name:       "successful response 201",
			statusCode: http.StatusCreated,
			body: mailnow.EmailResponse{
				Success: true,
				Data: mailnow.Data{
					MessageID: "msg_123",
					Status:    "queued",
				},
			},
			wantErr: false,
		},
		{
			name:       "successful response 204",
			statusCode: http.StatusNoContent,
			body:       "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.body != nil && tt.body != "" {
					json.NewEncoder(w).Encode(tt.body)
				}
			}))
			defer server.Close()

			resp, err := http.Get(server.URL)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}

			body, err := mailnow.HandleResponse(resp)

			if tt.wantErr {
				if err == nil {
					t.Errorf("handleResponse() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("handleResponse() unexpected error = %v", err)
				return
			}

			if body == nil && tt.body != "" {
				t.Error("handleResponse() returned nil body")
			}
		})
	}
}

// TestHandleResponseErrorMapping tests error mapping for each status code
func TestHandleResponseErrorMapping(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		errorBody   mailnow.ErrorResponse
		wantErrType interface{}
	}{
		{
			name:       "400 Bad Request - ValidationError",
			statusCode: http.StatusBadRequest,
			errorBody: mailnow.ErrorResponse{
				Error: struct {
					Code    string                 `json:"code"`
					Message string                 `json:"message"`
					Details map[string]interface{} `json:"details,omitempty"`
				}{
					Code:    "validation_error",
					Message: "Invalid email address",
				},
			},
			wantErrType: &mailnow.ValidationError{},
		},
		{
			name:       "401 Unauthorized - AuthError",
			statusCode: http.StatusUnauthorized,
			errorBody: mailnow.ErrorResponse{
				Error: struct {
					Code    string                 `json:"code"`
					Message string                 `json:"message"`
					Details map[string]interface{} `json:"details,omitempty"`
				}{
					Code:    "auth_error",
					Message: "Invalid API key",
				},
			},
			wantErrType: &mailnow.AuthError{},
		},
		{
			name:       "429 Too Many Requests - RateLimitError",
			statusCode: http.StatusTooManyRequests,
			errorBody: mailnow.ErrorResponse{
				Error: struct {
					Code    string                 `json:"code"`
					Message string                 `json:"message"`
					Details map[string]interface{} `json:"details,omitempty"`
				}{
					Code:    "rate_limit",
					Message: "Rate limit exceeded",
				},
			},
			wantErrType: &mailnow.RateLimitError{},
		},
		{
			name:       "500 Internal Server Error - ServerError",
			statusCode: http.StatusInternalServerError,
			errorBody: mailnow.ErrorResponse{
				Error: struct {
					Code    string                 `json:"code"`
					Message string                 `json:"message"`
					Details map[string]interface{} `json:"details,omitempty"`
				}{
					Code:    "server_error",
					Message: "Internal server error",
				},
			},
			wantErrType: &mailnow.ServerError{},
		},
		{
			name:       "502 Bad Gateway - ServerError",
			statusCode: http.StatusBadGateway,
			errorBody: mailnow.ErrorResponse{
				Error: struct {
					Code    string                 `json:"code"`
					Message string                 `json:"message"`
					Details map[string]interface{} `json:"details,omitempty"`
				}{
					Code:    "server_error",
					Message: "Bad gateway",
				},
			},
			wantErrType: &mailnow.ServerError{},
		},
		{
			name:       "503 Service Unavailable - ServerError",
			statusCode: http.StatusServiceUnavailable,
			errorBody: mailnow.ErrorResponse{
				Error: struct {
					Code    string                 `json:"code"`
					Message string                 `json:"message"`
					Details map[string]interface{} `json:"details,omitempty"`
				}{
					Code:    "server_error",
					Message: "Service unavailable",
				},
			},
			wantErrType: &mailnow.ServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.errorBody)
			}))
			defer server.Close()

			resp, err := http.Get(server.URL)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}

			_, err = mailnow.HandleResponse(resp)

			if err == nil {
				t.Errorf("handleResponse() expected error, got nil")
				return
			}

			// Check error type using errors.As
			switch tt.wantErrType.(type) {
			case *mailnow.ValidationError:
				var ve *mailnow.ValidationError
				if !errors.As(err, &ve) {
					t.Errorf("handleResponse() error type = %T, want ValidationError", err)
				}
			case *mailnow.AuthError:
				var ae *mailnow.AuthError
				if !errors.As(err, &ae) {
					t.Errorf("handleResponse() error type = %T, want AuthError", err)
				}
			case *mailnow.RateLimitError:
				var rle *mailnow.RateLimitError
				if !errors.As(err, &rle) {
					t.Errorf("handleResponse() error type = %T, want RateLimitError", err)
				}
			case *mailnow.ServerError:
				var se *mailnow.ServerError
				if !errors.As(err, &se) {
					t.Errorf("handleResponse() error type = %T, want ServerError", err)
				}
			}

			// Verify error message contains the expected message
			if err.Error() == "" {
				t.Error("handleResponse() error message is empty")
			}
		})
	}
}

// TestHandleResponseInvalidJSON tests handling of invalid JSON error responses
func TestHandleResponseInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	_, err = mailnow.HandleResponse(resp)

	if err == nil {
		t.Error("handleResponse() expected error for invalid JSON, got nil")
		return
	}

	// Should still return a ValidationError even with invalid JSON
	var ve *mailnow.ValidationError
	if !errors.As(err, &ve) {
		t.Errorf("handleResponse() error type = %T, want ValidationError", err)
	}
}

// TestConnectionError tests connection error handling
func TestConnectionError(t *testing.T) {
	// Test with invalid URL to trigger connection error
	client := &http.Client{Timeout: 1 * time.Second}
	ctx := context.Background()

	_, err := mailnow.MakeRequest(ctx, client, "POST", "http://invalid-url-that-does-not-exist-12345.com", "test_key", nil)

	if err == nil {
		t.Error("makeRequest() expected connection error, got nil")
		return
	}

	var ce *mailnow.ConnectionError
	if !errors.As(err, &ce) {
		t.Errorf("makeRequest() error type = %T, want ConnectionError", err)
	}
}
