# Design Document

## Overview

The Mailnow Go SDK is a lightweight, idiomatic Go client library that wraps the Mailnow email API. The design follows Go best practices with comprehensive error handling, context support, and Test-Driven Development (TDD) methodology. The SDK provides a simple, intuitive interface for sending emails using a struct-based architecture with clear separation between the client interface, HTTP communication, validation, and error handling.

## Architecture

The SDK follows a layered architecture:

```
┌─────────────────────────────────┐
│   Customer Application Code    │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│     Client (Public API)         │
│  - SendEmail()                  │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│   Validation & Serialization    │
│  - validateEmail()              │
│  - validateAPIKey()             │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│      HTTP Client Layer          │
│  - makeRequest()                │
│  - handleResponse()             │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│      Mailnow API Service        │
│  https://api.mailnow.xyz        │
└─────────────────────────────────┘
```

## Components and Interfaces

### 1. Client Struct

The main entry point for customers. This struct handles initialization and provides the public API.

**Location:** `client.go`

**Interface:**
```go
package mailnow

import "context"

// Client represents a Mailnow API client
type Client struct {
    apiKey     string
    httpClient *http.Client
    baseURL    string
}

// NewClient creates a new Mailnow client with the provided API key
func NewClient(apiKey string) (*Client, error)

// SendEmail sends an email via the Mailnow API
func (c *Client) SendEmail(ctx context.Context, req *EmailRequest) (*EmailResponse, error)
```

### 2. Request and Response Types

**Location:** `types.go`

```go
// EmailRequest represents an email sending request
type EmailRequest struct {
    From    string `json:"from"`
    To      string `json:"to"`
    Subject string `json:"subject"`
    HTML    string `json:"html"`
}

// EmailResponse represents a successful email sending response
type EmailResponse struct {
    Success   bool   `json:"success"`
    MessageID string `json:"message_id"`
    Status    string `json:"status"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
    Error struct {
        Code    string                 `json:"code"`
        Message string                 `json:"message"`
        Details map[string]interface{} `json:"details,omitempty"`
    } `json:"error"`
}
```

### 3. Validation Functions

Handles input validation before making API requests.

**Location:** `validation.go`

**Functions:**
```go
func validateAPIKey(apiKey string) error
func validateEmailAddress(email string) error
func validateEmailRequest(req *EmailRequest) error
```

### 4. HTTP Client Functions

Manages HTTP communication with the Mailnow API.

**Location:** `http.go`

**Functions:**
```go
func makeRequest(ctx context.Context, client *http.Client, method, url, apiKey string, body interface{}) (*http.Response, error)
func handleResponse(resp *http.Response) ([]byte, error)
```

### 5. Error Types

Custom error types for different error scenarios.

**Location:** `errors.go`

**Types:**
```go
// Error represents the base error type for all Mailnow SDK errors
type Error struct {
    Message string
    Err     error
}

func (e *Error) Error() string
func (e *Error) Unwrap() error

// ValidationError represents input validation failures
type ValidationError struct {
    *Error
}

// AuthError represents authentication failures
type AuthError struct {
    *Error
}

// RateLimitError represents rate limit exceeded errors
type RateLimitError struct {
    *Error
}

// ServerError represents server errors (5xx)
type ServerError struct {
    *Error
}

// ConnectionError represents network connection failures
type ConnectionError struct {
    *Error
}
```

## Error Handling

### Error Mapping Strategy

The SDK maps HTTP status codes to specific error types:

- **400 Bad Request** → `ValidationError`
- **401 Unauthorized** → `AuthError`
- **429 Too Many Requests** → `RateLimitError`
- **5xx Server Errors** → `ServerError`
- **Network/Connection Errors** → `ConnectionError`

### Go Error Handling Patterns

The SDK follows Go error handling idioms:
- All errors are returned as the last return value
- Errors can be unwrapped using `errors.Unwrap()`
- Errors can be checked using `errors.Is()` and `errors.As()`
- Error messages are descriptive and actionable

## Testing Strategy

### Test-Driven Development (TDD) Approach

The SDK follows strict TDD methodology:

1. **Write tests first** before implementing any functionality
2. **Run tests** to confirm they fail (red phase)
3. **Implement minimal code** to make tests pass (green phase)
4. **Refactor** while keeping tests green
5. **Repeat** for each feature

### Unit Tests

Test individual components in isolation using Go's testing package:

- **Validation Functions**: Test with empty, invalid, and valid inputs
- **Error Types**: Test error creation, type assertions, and unwrapping
- **Client Initialization**: Test with various API key formats
- **HTTP Functions**: Test request building, response parsing, error mapping using `httptest`

### Integration Tests

- Test complete flow with valid credentials
- Test authentication, validation, and error scenarios
- Test context cancellation

### Test Environment

- Use `mn_test_*` API keys for integration testing
- Mock HTTP responses using `httptest` package
- Achieve minimum 80% code coverage
- Use table-driven tests where appropriate

## Dependencies

### Standard Library Only

The SDK uses only Go standard library packages:
- `net/http` - HTTP client
- `context` - Context support
- `encoding/json` - JSON encoding/decoding
- `errors` - Error handling
- `fmt` - String formatting
- `regexp` - Email validation
- `time` - Timeout configuration

## Package Structure

```
mailnow/
├── client.go            # Main Client struct
├── types.go             # Request/Response types
├── errors.go            # Custom error types
├── validation.go        # Input validation
├── http.go              # HTTP communication
├── constants.go         # Constants
├── client_test.go       # Client tests
├── validation_test.go   # Validation tests
├── http_test.go         # HTTP tests
├── errors_test.go       # Error tests
├── README.md            # Documentation
├── go.mod               # Go module file
├── go.sum               # Dependency checksums
└── LICENSE              # License file
```

## Configuration

### Constants

**Location:** `constants.go`

```go
const (
    APIBaseURL = "https://api.mailnow.xyz"
    APIVersion = "v1"
    EmailSendEndpoint = "/v1/email/send"
    RequestTimeout = 30 * time.Second
    APIKeyPrefixLive = "mn_live_"
    APIKeyPrefixTest = "mn_test_"
)
```

## Usage Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/yourusername/mailnow"
)

func main() {
    // Initialize client
    client, err := mailnow.NewClient("mn_live_7e59df7ce4a14545b443837804ec9722")
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Prepare email request
    req := &mailnow.EmailRequest{
        From:    "sender@example.com",
        To:      "recipient@example.com",
        Subject: "Test Email",
        HTML:    "<h1>Hello World</h1><p>This is a test email.</p>",
    }
    
    // Send email
    resp, err := client.SendEmail(ctx, req)
    if err != nil {
        switch e := err.(type) {
        case *mailnow.ValidationError:
            log.Printf("Validation error: %v", e)
        case *mailnow.AuthError:
            log.Printf("Authentication error: %v", e)
        default:
            log.Printf("Error: %v", e)
        }
        return
    }
    
    fmt.Printf("Email sent: %s (status: %s)\n", resp.MessageID, resp.Status)
}
```

## Design Decisions

### 1. Standard Library Only

**Decision:** Use only Go standard library packages, no external dependencies.

**Rationale:**
- Reduces dependency management complexity
- Improves security by minimizing attack surface
- Ensures long-term stability
- Simplifies installation and updates

### 2. Context Support

**Decision:** Require context.Context for all API calls.

**Rationale:**
- Follows Go best practices
- Enables timeout control per request
- Supports cancellation for long-running operations
- Required for production-grade Go libraries

### 3. Struct-Based Configuration

**Decision:** Use struct for email request parameters.

**Rationale:**
- More extensible for future additions
- Cleaner function signatures
- Common Go pattern for complex parameters

### 4. Error Types vs Error Values

**Decision:** Use custom error types (structs) instead of sentinel error values.

**Rationale:**
- Allows type assertions for error handling
- Can include additional context
- Supports error wrapping
- More flexible than sentinel errors

### 5. Exported vs Unexported

**Decision:** Export only Client, request/response types, and error types.

**Rationale:**
- Minimal public API
- Reduces API surface area
- Allows internal refactoring without breaking changes

## Security Considerations

1. **API Key Storage**: API keys stored in memory only, never logged
2. **HTTPS Only**: All API communication uses HTTPS
3. **Input Validation**: All parameters validated before sending
4. **Error Messages**: Avoid exposing sensitive information
5. **No Credential Logging**: SDK never logs API keys

## Go Version Support

**Decision:** Support Go 1.21 and above.

**Rationale:**
- Modern Go features
- Aligns with Go's support policy
- Most production environments use recent versions

## Documentation Strategy

### 1. README.md

- Installation instructions via go get
- Getting started guide
- Code examples
- Error handling patterns

### 2. Godoc Comments

- Package-level documentation
- Type and function documentation
- Usage examples

### 3. Error Documentation

- When each error is returned
- How to handle each error type

## Package Distribution

### Module Configuration

**Location:** `go.mod`

- Module path: `github.com/yourusername/mailnow`
- Go version: 1.21
- No external dependencies

### Installation

```bash
go get github.com/yourusername/mailnow
```

## HTTP Connection Management

### Connection Reuse

- Client maintains single http.Client instance
- Connection pooling handled automatically by net/http
- Improves performance for multiple requests

### Timeout Configuration

- Default timeout: 30 seconds via http.Client.Timeout
- Can be overridden per-request using context.WithTimeout()

### Resource Cleanup

- http.Client automatically manages connections
- No explicit cleanup needed in most cases

## CI/CD Integration

### GitHub Actions Workflow

**Location:** `.github/workflows/test.yml`

**Features:**
- Run tests on multiple Go versions
- Run on multiple platforms
- Check code formatting with `gofmt`
- Run `go vet` for static analysis
- Generate coverage reports
- Fail build if coverage drops below 80%
