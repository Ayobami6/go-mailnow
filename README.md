# go-mailnow

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue.svg)](https://golang.org/doc/go1.21)
[![Go Report Card](https://goreportcard.com/badge/github.com/Ayobami6/go-mailnow)](https://goreportcard.com/report/github.com/Ayobami6/go-mailnow)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A lightweight, idiomatic Go client library for the [Mailnow](https://mailnow.xyz) email API. Send emails programmatically with a simple, clean interface.

## Features

- 🚀 Simple and intuitive API
- 🔒 Secure API key authentication
- ✅ Comprehensive input validation
- 🎯 Context support for timeouts and cancellation
- 📦 Zero external dependencies (standard library only)
- 🧪 Well-tested with high code coverage
- 🔧 Idiomatic Go design

## Installation

Install the package using `go get`:

```bash
go get github.com/Ayobami6/go-mailnow
```

## Getting Started

### Prerequisites

You'll need a Mailnow API key to use this SDK. Get your API key from the [Mailnow Dashboard](https://mailnow.xyz/dashboard).

API keys come in two formats:
- **Live keys**: `mn_live_*` - for production use
- **Test keys**: `mn_test_*` - for development and testing

### Basic Usage

Here's a simple example to send an email:

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/Ayobami6/go-mailnow"
)

func main() {
    // Initialize the client with your API key
    client, err := mailnow.NewClient("mn_live_your_api_key_here")
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    
    // Create an email request
    req := &mailnow.EmailRequest{
        From:    "sender@example.com",
        To:      "recipient@example.com",
        Subject: "Hello from go-mailnow!",
        HTML:    "<h1>Welcome</h1><p>This is a test email sent via go-mailnow SDK.</p>",
    }
    
    // Send the email
    resp, err := client.SendEmail(context.Background(), req)
    if err != nil {
        log.Fatalf("Failed to send email: %v", err)
    }
    
    fmt.Printf("Email sent successfully!\n")
    fmt.Printf("Message ID: %s\n", resp.MessageID)
    fmt.Printf("Status: %s\n", resp.Status)
}
```

## Complete Example with Error Handling

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/Ayobami6/go-mailnow"
)

func main() {
    // Initialize client
    client, err := mailnow.NewClient("mn_live_your_api_key_here")
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Prepare email request
    req := &mailnow.EmailRequest{
        From:    "noreply@yourcompany.com",
        To:      "customer@example.com",
        Subject: "Welcome to Our Service",
        HTML:    "<h1>Welcome!</h1><p>Thank you for signing up.</p>",
    }
    
    // Send email with comprehensive error handling
    resp, err := client.SendEmail(ctx, req)
    if err != nil {
        // Handle specific error types
        switch e := err.(type) {
        case *mailnow.ValidationError:
            log.Printf("Validation error: %v", e)
            log.Println("Please check your email parameters")
        case *mailnow.AuthError:
            log.Printf("Authentication error: %v", e)
            log.Println("Please check your API key")
        case *mailnow.RateLimitError:
            log.Printf("Rate limit exceeded: %v", e)
            log.Println("Please wait before sending more emails")
        case *mailnow.ServerError:
            log.Printf("Server error: %v", e)
            log.Println("Please try again later")
        case *mailnow.ConnectionError:
            log.Printf("Connection error: %v", e)
            log.Println("Please check your network connection")
        default:
            log.Printf("Unexpected error: %v", err)
        }
        return
    }
    
    // Success!
    fmt.Printf("✓ Email sent successfully!\n")
    fmt.Printf("  Message ID: %s\n", resp.MessageID)
    fmt.Printf("  Status: %s\n", resp.Status)
}
```

## Context Usage

### Timeout Example

Control request timeouts using context:

```go
// Set a 5-second timeout for the email request
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.SendEmail(ctx, req)
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        log.Println("Request timed out")
    }
    log.Fatal(err)
}
```

### Cancellation Example

Cancel requests programmatically:

```go
ctx, cancel := context.WithCancel(context.Background())

// Start sending email in a goroutine
go func() {
    resp, err := client.SendEmail(ctx, req)
    if err != nil {
        if ctx.Err() == context.Canceled {
            log.Println("Request was cancelled")
            return
        }
        log.Printf("Error: %v", err)
        return
    }
    log.Printf("Email sent: %s", resp.MessageID)
}()

// Cancel the request after some condition
time.Sleep(1 * time.Second)
cancel()
```

## Error Types

The SDK provides specific error types for different failure scenarios:

### ValidationError

Returned when input validation fails before making an API request.

**Common causes:**
- Empty or missing required fields (from, to, subject, html)
- Invalid email address format
- Invalid API key format

**Example:**
```go
resp, err := client.SendEmail(ctx, &mailnow.EmailRequest{
    From: "",  // Empty from field
    To: "recipient@example.com",
    Subject: "Test",
    HTML: "<p>Test</p>",
})
if err != nil {
    if _, ok := err.(*mailnow.ValidationError); ok {
        log.Println("Validation failed - check your input")
    }
}
```

### AuthError

Returned when authentication fails (HTTP 401).

**Common causes:**
- Invalid API key
- Expired API key
- API key not authorized for the requested operation

**Example:**
```go
client, err := mailnow.NewClient("mn_live_invalid_key")
if err != nil {
    if _, ok := err.(*mailnow.ValidationError); ok {
        log.Println("Invalid API key format")
    }
}

resp, err := client.SendEmail(ctx, req)
if err != nil {
    if _, ok := err.(*mailnow.AuthError); ok {
        log.Println("Authentication failed - check your API key")
    }
}
```

### RateLimitError

Returned when you exceed the API rate limits (HTTP 429).

**Common causes:**
- Sending too many emails in a short time period
- Exceeding your plan's rate limits

**Example:**
```go
resp, err := client.SendEmail(ctx, req)
if err != nil {
    if _, ok := err.(*mailnow.RateLimitError); ok {
        log.Println("Rate limit exceeded - please wait before retrying")
        // Implement exponential backoff or wait before retrying
    }
}
```

### ServerError

Returned when the Mailnow API experiences server errors (HTTP 5xx).

**Common causes:**
- Temporary server issues
- API maintenance
- Internal server errors

**Example:**
```go
resp, err := client.SendEmail(ctx, req)
if err != nil {
    if _, ok := err.(*mailnow.ServerError); ok {
        log.Println("Server error - please try again later")
        // Implement retry logic with exponential backoff
    }
}
```

### ConnectionError

Returned when network connectivity issues occur.

**Common causes:**
- Network timeout
- DNS resolution failure
- Connection refused
- TLS/SSL errors

**Example:**
```go
resp, err := client.SendEmail(ctx, req)
if err != nil {
    if _, ok := err.(*mailnow.ConnectionError); ok {
        log.Println("Connection error - check your network")
    }
}
```

## API Reference

### Client

#### NewClient

```go
func NewClient(apiKey string) (*Client, error)
```

Creates a new Mailnow client with the provided API key.

**Parameters:**
- `apiKey` (string): Your Mailnow API key (format: `mn_live_*` or `mn_test_*`)

**Returns:**
- `*Client`: Initialized client instance
- `error`: ValidationError if the API key is invalid or empty

**Example:**
```go
client, err := mailnow.NewClient("mn_live_your_api_key_here")
if err != nil {
    log.Fatal(err)
}
```

#### SendEmail

```go
func (c *Client) SendEmail(ctx context.Context, req *EmailRequest) (*EmailResponse, error)
```

Sends an email via the Mailnow API.

**Parameters:**
- `ctx` (context.Context): Context for timeout and cancellation control
- `req` (*EmailRequest): Email request containing from, to, subject, and html fields

**Returns:**
- `*EmailResponse`: Response containing message ID and status
- `error`: One of ValidationError, AuthError, RateLimitError, ServerError, or ConnectionError

**Example:**
```go
req := &mailnow.EmailRequest{
    From:    "sender@example.com",
    To:      "recipient@example.com",
    Subject: "Hello",
    HTML:    "<p>Hello World</p>",
}

resp, err := client.SendEmail(context.Background(), req)
if err != nil {
    log.Fatal(err)
}
```

### Types

#### EmailRequest

```go
type EmailRequest struct {
    From    string `json:"from"`
    To      string `json:"to"`
    Subject string `json:"subject"`
    HTML    string `json:"html"`
}
```

Represents an email sending request.

**Fields:**
- `From` (string): Sender email address (required)
- `To` (string): Recipient email address (required)
- `Subject` (string): Email subject line (required)
- `HTML` (string): Email body in HTML format (required)

#### EmailResponse

```go
type EmailResponse struct {
    Success   bool   `json:"success"`
    MessageID string `json:"message_id"`
    Status    string `json:"status"`
}
```

Represents a successful email sending response.

**Fields:**
- `Success` (bool): Whether the email was sent successfully
- `MessageID` (string): Unique identifier for the sent email
- `Status` (string): Current status of the email

## Requirements

- Go 1.21 or higher
- No external dependencies (uses standard library only)

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Resources

- [Mailnow API Documentation](https://docs.mailnow.xyz)
- [Mailnow Dashboard](https://mailnow.xyz/dashboard)
- [Go Documentation](https://pkg.go.dev/github.com/Ayobami6/go-mailnow)

## Support

If you encounter any issues or have questions:
- Check the [API Documentation](https://docs.mailnow.xyz)
- Open an issue on [GitHub](https://github.com/Ayobami6/go-mailnow/issues)
- Contact Mailnow support at support@mailnow.xyz

## Changelog

### v1.0.0
- Initial release
- Basic email sending functionality
- Comprehensive error handling
- Context support for timeouts and cancellation
- Full test coverage