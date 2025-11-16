# Requirements Document

## Introduction

The Mailnow Go SDK is a client library that enables customers of the Mailnow SaaS email service to send emails programmatically from their Go applications. The SDK provides a simple, idiomatic Go interface to interact with the Mailnow public email API (https://api.mailnow.xyz) using customer API keys for authentication.

## Glossary

- **Mailnow SDK**: The Go client library that wraps the Mailnow email API
- **Mailnow API**: The RESTful API service hosted at https://api.mailnow.xyz that processes email sending requests
- **API Key**: A customer-specific authentication token (format: mn_live_* or mn_test_*) used to authenticate requests to the Mailnow API
- **Email Client**: The main SDK struct that customers instantiate to send emails
- **Customer**: A user of the Mailnow SaaS service who uses the Go SDK in their application

## Requirements

### Requirement 1

**User Story:** As a customer, I want to initialize the Mailnow SDK with my API key, so that I can authenticate my email sending requests.

#### Acceptance Criteria

1. THE Mailnow SDK SHALL provide a client constructor function that accepts an API key as a parameter
2. WHEN a customer instantiates the client with an empty API key, THE Mailnow SDK SHALL return a validation error
3. THE Mailnow SDK SHALL store the API key securely within the client instance for subsequent API requests
4. WHERE the API key format is invalid (does not start with "mn_live_" or "mn_test_"), THE Mailnow SDK SHALL return a validation error

### Requirement 2

**User Story:** As a customer, I want to send an email with basic fields (from, to, subject, HTML body), so that I can deliver messages to recipients.

#### Acceptance Criteria

1. THE Mailnow SDK SHALL provide a SendEmail method that accepts from, to, subject, and html parameters
2. WHEN the SendEmail method is called, THE Mailnow SDK SHALL make an HTTP POST request to https://api.mailnow.xyz/v1/email/send
3. THE Mailnow SDK SHALL include the API key in the X-API-Key header for authentication
4. THE Mailnow SDK SHALL include Content-Type: application/json header in the request
5. THE Mailnow SDK SHALL serialize the email parameters as JSON in the request body

### Requirement 3

**User Story:** As a customer, I want to receive clear feedback when my email is sent successfully, so that I can confirm the operation completed.

#### Acceptance Criteria

1. WHEN the Mailnow API returns a successful response (HTTP 2xx), THE Mailnow SDK SHALL return a success response to the customer
2. THE Mailnow SDK SHALL parse the API response and return relevant data (such as message ID or status)
3. THE Mailnow SDK SHALL provide the response in a structured format (Go struct)

### Requirement 4

**User Story:** As a customer, I want to receive clear error messages when email sending fails, so that I can understand and fix the issue.

#### Acceptance Criteria

1. WHEN the Mailnow API returns an authentication error (HTTP 401), THE Mailnow SDK SHALL return an authentication error with a descriptive message
2. WHEN the Mailnow API returns a validation error (HTTP 400), THE Mailnow SDK SHALL return a validation error with details about invalid fields
3. WHEN the Mailnow API returns a rate limit error (HTTP 429), THE Mailnow SDK SHALL return a rate limit error
4. WHEN the Mailnow API returns a server error (HTTP 5xx), THE Mailnow SDK SHALL return a server error
5. WHEN a network error occurs, THE Mailnow SDK SHALL return a connection error with relevant error details

### Requirement 5

**User Story:** As a customer, I want to validate email parameters before sending, so that I can catch errors early in my application.

#### Acceptance Criteria

1. WHEN the from parameter is missing or empty, THE Mailnow SDK SHALL return a validation error before making the API request
2. WHEN the to parameter is missing or empty, THE Mailnow SDK SHALL return a validation error before making the API request
3. WHEN the subject parameter is missing or empty, THE Mailnow SDK SHALL return a validation error before making the API request
4. WHEN the html parameter is missing or empty, THE Mailnow SDK SHALL return a validation error before making the API request
5. WHEN email addresses are provided in an invalid format, THE Mailnow SDK SHALL return a validation error with details about the invalid address

### Requirement 6

**User Story:** As a customer, I want to install the SDK easily via go get, so that I can quickly integrate it into my Go projects.

#### Acceptance Criteria

1. THE Mailnow SDK SHALL be installable via go get command
2. THE Mailnow SDK SHALL use Go modules for dependency management
3. THE Mailnow SDK SHALL support Go 1.21 and above
4. THE Mailnow SDK SHALL include a README with installation and basic usage instructions

### Requirement 7

**User Story:** As a customer, I want clear documentation and examples, so that I can quickly understand how to use the SDK.

#### Acceptance Criteria

1. THE Mailnow SDK SHALL include a README file with installation instructions
2. THE Mailnow SDK SHALL include code examples demonstrating basic email sending
3. THE Mailnow SDK SHALL include godoc comments for all exported types and functions
4. THE Mailnow SDK SHALL document all error types that can be returned

### Requirement 8

**User Story:** As a customer, I want the SDK to handle HTTP connections efficiently, so that my application performs well.

#### Acceptance Criteria

1. THE Mailnow SDK SHALL use the standard net/http package for HTTP communication
2. THE Mailnow SDK SHALL set appropriate timeout values for API requests
3. THE Mailnow SDK SHALL properly manage HTTP connections using connection pooling
4. WHERE multiple emails are sent, THE Mailnow SDK SHALL reuse HTTP connections when possible

### Requirement 9

**User Story:** As a customer, I want the SDK to follow Go idioms and best practices, so that it integrates naturally with my Go codebase.

#### Acceptance Criteria

1. THE Mailnow SDK SHALL return errors as the last return value following Go conventions
2. THE Mailnow SDK SHALL use context.Context for cancellation and timeout control
3. THE Mailnow SDK SHALL export only necessary types and functions
4. THE Mailnow SDK SHALL follow standard Go naming conventions
5. THE Mailnow SDK SHALL pass go vet and golint checks without warnings
